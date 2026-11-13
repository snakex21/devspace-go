//go:build windows

package server

import (
	"errors"

	"golang.org/x/sys/windows"
)

func isAddressInUse(err error) bool {
	// WSAEACCES is returned when another socket holds the port exclusively.
	return errors.Is(err, windows.WSAEADDRINUSE) || errors.Is(err, windows.WSAEACCES)
}
