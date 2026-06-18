package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsPathInsideRoot checks if a path is inside the given root directory.
func IsPathInsideRoot(targetPath, root string) bool {
	// Clean and resolve paths
	cleanPath := filepath.Clean(targetPath)
	cleanRoot := filepath.Clean(root)

	// Handle Windows drive letters
	cleanPath = strings.ToLower(cleanPath)
	cleanRoot = strings.ToLower(cleanRoot)

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
	skipDirs := map[string]bool{
		".git":         true,
		".hg":          true,
		".svn":         true,
		".devspace":       true,
		".devspace-state": true,
		".webcoder":       true,
		".webcoder-state": true,
		"node_modules": true,
		"dist":         true,
		"build":        true,
		".next":        true,
		".turbo":       true,
		".cache":       true,
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip inaccessible paths
		}

		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		return visitor(path, info)
	})
}
