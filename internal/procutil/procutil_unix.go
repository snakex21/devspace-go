//go:build unix

package procutil

import (
	"os/exec"
	"syscall"
)

type osTree struct {
	pgid int
}

func configure(cmd *exec.Cmd, killWithParent bool) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	setParentDeathSignal(cmd.SysProcAttr, killWithParent)
}

func (t *Tree) attach() {
	t.pgid = t.cmd.Process.Pid
}

func (t *Tree) killTree() bool {
	if t.pgid <= 0 {
		return false
	}
	return syscall.Kill(-t.pgid, syscall.SIGKILL) == nil
}

func (t *Tree) release() {}
