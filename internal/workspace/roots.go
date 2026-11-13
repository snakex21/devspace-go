package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// The instruction-file scan run by open_workspace is bounded on every
// platform so that opening a huge root (for example a whole drive) returns
// quickly instead of walking millions of files.
const (
	workspaceWalkMaxFiles = 5000
	workspaceWalkTimeout  = 3 * time.Second
)

var workspaceSkippedDirs = map[string]bool{
	".git": true, ".hg": true, ".svn": true,
	".devspace": true, ".devspace-state": true,
	".webcoder": true, ".webcoder-state": true,
	"node_modules": true, "dist": true, "build": true,
	".next": true, ".turbo": true, ".cache": true,
}

var unixWorkspaceSkippedDirs = map[string]bool{
	".android":  true,
	".bun":      true,
	".config":   true,
	".gradle":   true,
	".hermes":   true,
	".local":    true,
	".npm":      true,
	".opencode": true,
	".var":      true,
}

// IsPathInsideRoot checks if a path is inside the given root directory.
func IsPathInsideRoot(targetPath, root string) bool {
	// Clean and resolve paths
	cleanPath := filepath.Clean(targetPath)
	cleanRoot := filepath.Clean(root)

	// Windows and default macOS volumes are case-insensitive. Linux is not:
	// folding there would treat /srv/App as if it were inside /srv/app.
	if pathsAreCaseInsensitive {
		cleanPath = strings.ToLower(cleanPath)
		cleanRoot = strings.ToLower(cleanRoot)
	}

	if !strings.HasSuffix(cleanRoot, string(filepath.Separator)) {
		cleanRoot += string(filepath.Separator)
	}

	return strings.HasPrefix(cleanPath, cleanRoot) || cleanPath == strings.TrimSuffix(cleanRoot, string(filepath.Separator))
}

// ResolvePath resolves a relative or absolute path against a working directory and validates
// it is inside one of the allowed roots.
func ResolvePath(inputPath, cwd string, allowedRoots []string) (string, error) {
	// If the path is absolute, use it directly
	var resolved string
	if filepath.IsAbs(inputPath) {
		resolved = filepath.Clean(inputPath)
	} else {
		// Handle ~/ paths
		if strings.HasPrefix(inputPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("cannot resolve home directory: %w", err)
			}
			resolved = filepath.Clean(filepath.Join(homeDir, inputPath[2:]))
		} else {
			resolved = filepath.Clean(filepath.Join(cwd, inputPath))
		}
	}

	// Check if resolved path is inside any allowed root
	for _, root := range allowedRoots {
		cleanRoot := filepath.Clean(root)
		if IsPathInsideRoot(resolved, cleanRoot) {
			return resolved, nil
		}
	}

	return "", fmt.Errorf("path %s is outside allowed roots", inputPath)
}

// ResolveInRoot resolves inputPath (relative to root, or absolute) and verifies
// that the result stays inside root, both lexically and after following any
// symlinks or junctions in the part of the path that already exists. An empty
// input resolves to root itself. The returned path is the cleaned, unresolved
// path so messages and relative formatting stay stable.
//
// A link created after this check (a TOCTOU race by a local process) is out of
// scope: such a process can already write anywhere the server can.
func ResolveInRoot(root, inputPath string) (string, error) {
	cleanRoot := filepath.Clean(root)
	if !filepath.IsAbs(cleanRoot) {
		absRoot, err := filepath.Abs(cleanRoot)
		if err != nil {
			return "", fmt.Errorf("resolve workspace root: %w", err)
		}
		cleanRoot = absRoot
	}
	if strings.ContainsRune(inputPath, 0) {
		return "", fmt.Errorf("path contains a NUL byte")
	}

	var candidate string
	switch {
	case strings.TrimSpace(inputPath) == "":
		candidate = cleanRoot
	case filepath.IsAbs(inputPath):
		candidate = filepath.Clean(inputPath)
	case strings.HasPrefix(inputPath, "~/") || strings.HasPrefix(inputPath, `~\`):
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot resolve home directory: %w", err)
		}
		candidate = filepath.Clean(filepath.Join(homeDir, inputPath[2:]))
	default:
		// A Windows drive-relative or rooted path ("C:x", `\x`) is neither
		// absolute nor relative to root; joining it would hide that ambiguity.
		if filepath.VolumeName(inputPath) != "" || (runtime.GOOS == "windows" && (strings.HasPrefix(inputPath, `\`) || strings.HasPrefix(inputPath, "/"))) {
			return "", fmt.Errorf("path %s is outside workspace root", inputPath)
		}
		candidate = filepath.Clean(filepath.Join(cleanRoot, inputPath))
	}

	if !IsPathInsideRoot(candidate, cleanRoot) {
		return "", fmt.Errorf("path %s is outside workspace root", inputPath)
	}

	// Fast path: if no existing component below the root is a link or other
	// reparse point, the path cannot leave the root, whatever the root itself
	// points to. This costs one Lstat per component instead of two full
	// EvalSymlinks walks from the volume root.
	hasLink, err := hasLinkBelowRoot(cleanRoot, candidate)
	if err != nil {
		return "", fmt.Errorf("resolve path %s: %w", inputPath, err)
	}
	if !hasLink {
		return candidate, nil
	}

	realRoot, err := filepath.EvalSymlinks(cleanRoot)
	if err != nil {
		// The root may not exist yet (open_workspace creates it); nothing below
		// it can be a link then.
		return candidate, nil
	}
	realCandidate, err := evalExistingPrefix(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve path %s: %w", inputPath, err)
	}
	if !IsPathInsideRoot(realCandidate, realRoot) {
		return "", fmt.Errorf("path %s resolves outside workspace root through a link", inputPath)
	}
	return candidate, nil
}

// hasLinkBelowRoot reports whether any existing component of candidate below
// root is a symlink, junction or other reparse point (Go reports Windows
// junctions and non-symlink reparse points as ModeIrregular). Components after
// the first missing one cannot be links. Anything it cannot classify is
// reported as a link, so the caller falls back to the full check.
func hasLinkBelowRoot(root, candidate string) (bool, error) {
	rel, err := filepath.Rel(root, candidate)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return true, nil
	}
	if rel == "." {
		return false, nil
	}
	current := root
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			if os.IsNotExist(err) || isNotDirectory(err) {
				return false, nil
			}
			return false, err
		}
		if info.Mode()&(os.ModeSymlink|os.ModeIrregular) != 0 {
			return true, nil
		}
	}
	return false, nil
}

// evalExistingPrefix resolves links in the longest existing prefix of path and
// re-appends the components that do not exist yet.
func evalExistingPrefix(path string) (string, error) {
	existing := path
	var missing []string
	for {
		if _, err := os.Lstat(existing); err == nil {
			break
		} else if !os.IsNotExist(err) && !isNotDirectory(err) {
			return "", err
		}
		parent := filepath.Dir(existing)
		if parent == existing {
			return path, nil
		}
		missing = append(missing, filepath.Base(existing))
		existing = parent
	}
	resolved, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", err
	}
	for i := len(missing) - 1; i >= 0; i-- {
		resolved = filepath.Join(resolved, missing[i])
	}
	return resolved, nil
}

// AssertAllowedPath checks that the absolute path is inside one of the allowed roots.
func AssertAllowedPath(absPath string, allowedRoots []string) (string, error) {
	cleanPath := filepath.Clean(absPath)
	for _, root := range allowedRoots {
		cleanRoot := filepath.Clean(root)
		if IsPathInsideRoot(cleanPath, cleanRoot) {
			return cleanPath, nil
		}
	}
	return "", fmt.Errorf("path %s is outside allowed roots", absPath)
}

// WalkWorkspace walks a directory tree, skipping blacklisted directories.
func WalkWorkspace(root string, visitor func(path string, info os.FileInfo) error) error {
	if runtime.GOOS == "windows" {
		return walkWorkspaceWindows(root, visitor, workspaceWalkMaxFiles, workspaceWalkTimeout, time.Now)
	}
	return walkWorkspaceUnix(root, visitor, workspaceWalkMaxFiles, workspaceWalkTimeout, time.Now)
}

// walkWorkspaceWindows keeps the original Windows directory set (only
// workspaceSkippedDirs are pruned; dot-directories such as .hermes stay
// visible) but applies the same file-count and time limits as the Unix walk.
func walkWorkspaceWindows(
	root string,
	visitor func(path string, info os.FileInfo) error,
	maxFiles int,
	timeout time.Duration,
	now func() time.Time,
) error {
	started := now()
	filesSeen := 0

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip inaccessible paths
		}
		if timeout > 0 && now().Sub(started) >= timeout {
			return filepath.SkipAll
		}

		if info.IsDir() {
			if workspaceSkippedDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		filesSeen++
		if maxFiles > 0 && filesSeen > maxFiles {
			return filepath.SkipAll
		}

		return visitor(path, info)
	})
}

func walkWorkspaceUnix(
	root string,
	visitor func(path string, info os.FileInfo) error,
	maxFiles int,
	timeout time.Duration,
	now func() time.Time,
) error {
	cleanRoot := filepath.Clean(root)
	started := now()
	filesSeen := 0

	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip inaccessible paths
		}
		if timeout > 0 && now().Sub(started) >= timeout {
			return filepath.SkipAll
		}

		if entry.IsDir() {
			if filepath.Clean(path) != cleanRoot &&
				(workspaceSkippedDirs[entry.Name()] || unixWorkspaceSkippedDirs[entry.Name()]) {
				return filepath.SkipDir
			}
			return nil
		}

		filesSeen++
		if maxFiles > 0 && filesSeen > maxFiles {
			return filepath.SkipAll
		}

		info, err := entry.Info()
		if err != nil {
			return nil
		}
		return visitor(path, info)
	})
}
