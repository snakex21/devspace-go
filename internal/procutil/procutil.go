// Package procutil starts child processes so that their whole process tree
// can be terminated, and optionally dies together with this process.
//
// Windows: every child is placed in its own job object. Kill terminates the
// job (the child and all of its descendants). With KillWithParent the job is
// created with JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE, so the operating system
// kills the tree when Dev Space Go exits for any reason, including a crash or
// a hard kill from Task Manager.
//
// Unix: every child gets its own process group and Kill signals the group.
// With KillWithParent, Linux additionally sets PR_SET_PDEATHSIG. macOS has no
// equivalent, so an orphan can survive a hard kill (SIGKILL) there.
package procutil

import (
	"errors"
	"os"
	"os/exec"
	"sync"
)

// Tree owns a started command and its descendants.
type Tree struct {
	cmd            *exec.Cmd
	killWithParent bool

	mu       sync.Mutex
	released bool
	osTree
}

// Prepare configures cmd for tree management. It must be called before
// Start. If cmd was created with exec.CommandContext, cancelling the context
// kills the whole tree instead of only the direct child.
func Prepare(cmd *exec.Cmd, killWithParent bool) *Tree {
	t := &Tree{cmd: cmd, killWithParent: killWithParent}
	configure(cmd, killWithParent)
	// exec.CommandContext installs a default Cancel that kills only the
	// direct child; replace it. Commands without a context keep Cancel nil.
	if cmd.Cancel != nil {
		cmd.Cancel = func() error {
			err := t.Kill()
			if errors.Is(err, os.ErrProcessDone) {
				return nil
			}
			return err
		}
	}
	return t
}

// Start starts the command and attaches it to its tree. If attaching fails
// the command still runs; Kill then falls back to killing the direct child.
func (t *Tree) Start() error {
	if err := t.cmd.Start(); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.attach()
	return nil
}

// Kill terminates the command and all processes it started.
func (t *Tree) Kill() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.cmd.Process == nil {
		return nil
	}
	if t.killTree() {
		return nil
	}
	return t.cmd.Process.Kill()
}

// Release frees tree resources after Wait has returned. For a tree created
// with killWithParent, descendants that are still running are killed.
func (t *Tree) Release() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.released {
		return
	}
	t.released = true
	t.release()
}
