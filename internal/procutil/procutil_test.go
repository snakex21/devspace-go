package procutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

// TestHelperProcess is re-executed as a child process by the tests below.
//
//	heartbeat FILE           append to FILE every 20ms for at most 30s
//	spawn FILE               start "heartbeat FILE" as a child, then sleep
//	spawn-managed FILE       like spawn, but the child is started through
//	                         Prepare(killWithParent=true)
func TestHelperProcess(t *testing.T) {
	if os.Getenv("DEVSPACE_PROCUTIL_HELPER") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 && args[0] != "--" {
		args = args[1:]
	}
	mode, file := args[1], args[2]
	switch mode {
	case "heartbeat":
		os.WriteFile(file+".pid", []byte(strconv.Itoa(os.Getpid())), 0o600)
		deadline := time.Now().Add(30 * time.Second)
		for time.Now().Before(deadline) {
			f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
			if err == nil {
				f.WriteString(".")
				f.Close()
			}
			time.Sleep(20 * time.Millisecond)
		}
		os.Exit(0)
	case "spawn", "spawn-managed":
		child := helper(context.Background(), "heartbeat", file)
		if mode == "spawn-managed" {
			if err := Prepare(child, true).Start(); err != nil {
				os.Exit(3)
			}
		} else if err := child.Start(); err != nil {
			os.Exit(3)
		}
		fmt.Println("started")
		time.Sleep(30 * time.Second)
		os.Exit(0)
	}
	os.Exit(2)
}

func helper(ctx context.Context, mode, file string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHelperProcess$", "--", mode, file)
	cmd.Env = append(os.Environ(), "DEVSPACE_PROCUTIL_HELPER=1")
	return cmd
}

func waitForHeartbeat(t *testing.T, file string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if info, err := os.Stat(file); err == nil && info.Size() > 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("grandchild never started")
}

// requireHeartbeatStops fails if file keeps growing, i.e. the writer survived.
func requireHeartbeatStops(t *testing.T, file string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		before, _ := os.ReadFile(file)
		time.Sleep(300 * time.Millisecond)
		after, _ := os.ReadFile(file)
		if len(before) == len(after) {
			return
		}
	}
	t.Fatal("grandchild is still running")
}

func TestContextCancelKillsWholeTree(t *testing.T) {
	file := filepath.Join(t.TempDir(), "heartbeat")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := helper(ctx, "spawn", file)
	tree := Prepare(cmd, false)
	if err := tree.Start(); err != nil {
		t.Fatal(err)
	}
	defer tree.Release()
	waitForHeartbeat(t, file)

	cancel()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Wait did not return after cancellation")
	}
	requireHeartbeatStops(t, file)
}

func TestKillWithParentKillsTreeOnParentHardKill(t *testing.T) {
	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
		t.Skip("no kernel support for killing children with their parent on " + runtime.GOOS)
	}
	file := filepath.Join(t.TempDir(), "heartbeat")
	// The intermediate process simulates Dev Space Go: it starts a managed
	// child (like cloudflared) and is then killed without any cleanup.
	parent := helper(context.Background(), "spawn-managed", file)
	if err := parent.Start(); err != nil {
		t.Fatal(err)
	}
	waitForHeartbeat(t, file)
	if err := parent.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	parent.Wait()
	requireHeartbeatStops(t, file)
}

func TestKillAfterExitIsHarmless(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	tree := Prepare(cmd, true)
	if err := tree.Start(); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	_ = tree.Kill()
	tree.Release()
	tree.Release()
}

func TestNormalExitLeavesBackgroundChildrenRunning(t *testing.T) {
	// Without killWithParent, releasing a finished tree must not kill
	// processes the command deliberately left running in the background.
	file := filepath.Join(t.TempDir(), "heartbeat")
	cmd := helper(context.Background(), "spawn", file)
	tree := Prepare(cmd, false)
	if err := tree.Start(); err != nil {
		t.Fatal(err)
	}
	waitForHeartbeat(t, file)
	// Kill only the direct child, simulating it exiting by itself.
	cmd.Process.Kill()
	cmd.Wait()
	tree.Release()
	before, _ := os.ReadFile(file)
	time.Sleep(300 * time.Millisecond)
	after, _ := os.ReadFile(file)
	if len(after) == len(before) {
		t.Fatal("background grandchild was killed on release")
	}
	killHeartbeat(t, file)
}

// killHeartbeat stops a leftover heartbeat helper by its recorded PID;
// heartbeats also exit on their own after 30s.
func killHeartbeat(t *testing.T, file string) {
	t.Helper()
	data, err := os.ReadFile(file + ".pid")
	if err != nil {
		t.Logf("heartbeat pid unknown: %v", err)
		return
	}
	pid, _ := strconv.Atoi(string(data))
	if process, err := os.FindProcess(pid); err == nil {
		process.Kill()
	}
}
