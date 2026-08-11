//go:build windows

package config

import (
	"strings"
	"syscall"
	"unsafe"
)

const localeNameMaxLength = 85

var getUserDefaultLocaleName = syscall.NewLazyDLL("kernel32.dll").NewProc("GetUserDefaultLocaleName")

func detectOSSystemLang() string {
	buffer := make([]uint16, localeNameMaxLength)
	length, _, _ := getUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	if length == 0 {
		return ""
	}

	locale := syscall.UTF16ToString(buffer)
	if len(locale) < 2 {
		return ""
	}
	return strings.ToLower(locale[:2])
}
