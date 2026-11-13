package workspace

import (
	"errors"
	"runtime"
	"syscall"
)

// pathsAreCaseInsensitive reports whether path containment should ignore case.
// Windows and default macOS (APFS/HFS+) volumes are case-insensitive.
var pathsAreCaseInsensitive = runtime.GOOS == "windows" || runtime.GOOS == "darwin"

func isNotDirectory(err error) bool {
	return errors.Is(err, syscall.ENOTDIR)
}
