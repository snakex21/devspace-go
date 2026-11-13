//go:build windows

package procutil

import (
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

type osTree struct {
	job windows.Handle
}

func configure(*exec.Cmd, bool) {}

// attach must be called with t.mu held, right after Start. A grandchild
// spawned before this point escapes the job; children normally need far
// longer than this to start their own children.
func (t *Tree) attach() {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return
	}
	if t.killWithParent {
		info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
		info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
		if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation,
			uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil {
			windows.CloseHandle(job)
			return
		}
	}
	process, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(t.cmd.Process.Pid))
	if err != nil {
		windows.CloseHandle(job)
		return
	}
	defer windows.CloseHandle(process)
	if err := windows.AssignProcessToJobObject(job, process); err != nil {
		windows.CloseHandle(job)
		return
	}
	t.job = job
}

func (t *Tree) killTree() bool {
	if t.job == 0 {
		return false
	}
	return windows.TerminateJobObject(t.job, 1) == nil
}

func (t *Tree) release() {
	if t.job != 0 {
		windows.CloseHandle(t.job)
		t.job = 0
	}
}
