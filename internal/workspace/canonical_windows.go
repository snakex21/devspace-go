//go:build windows

package workspace

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

// volumeNameDOS (with FILE_NAME_NORMALIZED, also 0) asks for a drive-letter path.
const volumeNameDOS = 0

// realPath asks Windows for the final path of an existing file or directory.
// Opening a handle follows symlinks, junctions and mount points (unlike
// filepath.EvalSymlinks, which since Go 1.23 leaves mount points and
// junctions unresolved), and GetFinalPathNameByHandle also expands 8.3 short
// names, normalises letter case and maps subst drives to their target.
func realPath(path string) (string, error) {
	pathp, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}
	handle, err := windows.CreateFile(pathp, 0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", path, err)
	}
	defer windows.CloseHandle(handle)

	buf := make([]uint16, windows.MAX_PATH)
	for {
		n, err := windows.GetFinalPathNameByHandle(handle, &buf[0], uint32(len(buf)), volumeNameDOS)
		if err != nil {
			return "", fmt.Errorf("final path of %s: %w", path, err)
		}
		if n < uint32(len(buf)) {
			return stripExtendedPrefix(windows.UTF16ToString(buf[:n]))
		}
		buf = make([]uint16, n+1)
	}
}

// stripExtendedPrefix turns \\?\C:\x into C:\x and \\?\UNC\srv\share into
// \\srv\share. Anything else (a volume GUID path of a volume without a drive
// letter) cannot be compared with configured roots and is refused.
func stripExtendedPrefix(path string) (string, error) {
	const extended, extendedUNC = `\\?\`, `\\?\UNC\`
	switch {
	case strings.HasPrefix(path, extendedUNC):
		return `\\` + path[len(extendedUNC):], nil
	case strings.HasPrefix(path, extended) && len(path) >= len(extended)+2 && path[len(extended)+1] == ':':
		return path[len(extended):], nil
	case strings.HasPrefix(path, extended):
		return "", fmt.Errorf("path %s is on a volume without a drive letter", path)
	}
	return path, nil
}
