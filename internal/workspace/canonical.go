package workspace

import (
	"os"
	"path/filepath"
)

// Canonicalize returns the real path behind an absolute path: symlinks,
// junctions, mount points, subst drives, 8.3 short names and letter case are
// resolved for the longest prefix that exists (by the operating system, not
// by string manipulation). Components that do not exist yet are appended
// unchanged; a missing component cannot be a link.
//
// Containment checks must compare canonical paths: a lexical check alone
// lets a junction or symlink inside an allowed folder point anywhere.
func Canonicalize(path string) (string, error) {
	path = filepath.Clean(path)
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
			// Nothing exists, not even the volume; nothing to resolve.
			return path, nil
		}
		missing = append(missing, filepath.Base(existing))
		existing = parent
	}
	resolved, err := realPath(existing)
	if err != nil {
		return "", err
	}
	for i := len(missing) - 1; i >= 0; i-- {
		resolved = filepath.Join(resolved, missing[i])
	}
	return resolved, nil
}
