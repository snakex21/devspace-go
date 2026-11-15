//go:build !windows

package workspace

import "path/filepath"

// realPath resolves symlinks in an existing path. Bind mounts are not links
// and are deliberately not resolved: they are configured by an administrator.
func realPath(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
