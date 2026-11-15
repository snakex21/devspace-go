//go:build !windows

package workspace

import "testing"

func shortPathName(t *testing.T, path string) string {
	t.Skip("8.3 names are Windows-only")
	return path
}
