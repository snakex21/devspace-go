//go:build linux

package procutil

import "syscall"

func setParentDeathSignal(attr *syscall.SysProcAttr, enabled bool) {
	if enabled {
		attr.Pdeathsig = syscall.SIGKILL
	}
}
