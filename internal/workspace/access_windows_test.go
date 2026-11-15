//go:build windows

package workspace

import (
	"testing"

	"golang.org/x/sys/windows"
)

func shortPathName(t *testing.T, path string) string {
	t.Helper()
	long, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]uint16, windows.MAX_PATH)
	n, err := windows.GetShortPathName(long, &buf[0], uint32(len(buf)))
	if err != nil {
		t.Skipf("GetShortPathName: %v", err)
	}
	return windows.UTF16ToString(buf[:n])
}
