//go:build unix && !linux

package procutil

import "syscall"

func setParentDeathSignal(*syscall.SysProcAttr, bool) {}
