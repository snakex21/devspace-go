package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestShellHelperProcess is re-executed by the tests below in place of a shell.
func TestShellHelperProcess(t *testing.T) {
	if os.Getenv("DEVSPACE_SHELL_HELPER") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 && args[0] != "--" {
		args = args[1:]
	}
	mode := args[1]
	switch mode {
	case "heartbeat":
		file := args[2]
		os.WriteFile(file+".pid", []byte(strconv.Itoa(os.Getpid())), 0o600)
		deadline := time.Now().Add(30 * time.Second)
		for time.Now().Before(deadline) {
			if f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600); err == nil {
				f.WriteString(".")
				f.Close()
			}
			time.Sleep(20 * time.Millisecond)
		}
		os.Exit(0)
	case "background-exit", "background-hang":
		// Like "dev-server &": the child inherits stdout/stderr.
		child := exec.Command(os.Args[0], "-test.run=^TestShellHelperProcess$", "--", "heartbeat", args[2])
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr
		if err := child.Start(); err != nil {
			os.Exit(3)
		}
		fmt.Println("shell done")
		if mode == "background-hang" {
			time.Sleep(time.Minute)
		}
		os.Exit(0)
	case "flood":
		line := strings.Repeat("x", 1023) + "\n"
		for i := 0; i < 10*1024; i++ { // 10 MiB
			os.Stdout.WriteString(line)
		}
		os.Exit(0)
	}
	os.Exit(2)
}

func shellHelper(t *testing.T, mode string, extra ...string) (string, []string) {
	t.Helper()
	t.Setenv("DEVSPACE_SHELL_HELPER", "1")
	return os.Args[0], append([]string{"-test.run=^TestShellHelperProcess$", "--", mode}, extra...)
}

func killHelperByPIDFile(file string) {
	data, err := os.ReadFile(file + ".pid")
	if err != nil {
		return
	}
	pid, _ := strconv.Atoi(string(data))
	if process, err := os.FindProcess(pid); err == nil {
		process.Kill()
	}
}

func heartbeatGrows(file string) bool {
	before, _ := os.ReadFile(file)
	time.Sleep(300 * time.Millisecond)
	after, _ := os.ReadFile(file)
	return len(after) > len(before)
}

func TestRunCommandDoesNotHangOnBackgroundProcess(t *testing.T) {
	file := filepath.Join(t.TempDir(), "heartbeat")
	t.Cleanup(func() { killHelperByPIDFile(file) })
	name, args := shellHelper(t, "background-exit", file)

	started := time.Now()
	output, err := runCommand(context.Background(), t.TempDir(), name, args, 30)
	elapsed := time.Since(started)

	if elapsed > 10*time.Second {
		t.Fatalf("runCommand waited %v for a background process holding the output pipe", elapsed)
	}
	if err != nil {
		t.Fatalf("successful shell reported error: %v (output %q)", err, output)
	}
	if !strings.Contains(output, "shell done") || !strings.Contains(output, "background processes") {
		t.Fatalf("output = %q", output)
	}
	// The background process is intentionally left running.
	if !heartbeatGrows(file) {
		t.Fatal("background process was killed after a successful command")
	}
}

func TestRunCommandTimeoutKillsProcessTree(t *testing.T) {
	file := filepath.Join(t.TempDir(), "heartbeat")
	t.Cleanup(func() { killHelperByPIDFile(file) })
	name, args := shellHelper(t, "background-hang", file)

	started := time.Now()
	output, err := runCommand(context.Background(), t.TempDir(), name, args, 2)
	if elapsed := time.Since(started); elapsed > 10*time.Second {
		t.Fatalf("timeout took %v", elapsed)
	}
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("err = %v, output %q", err, output)
	}
	deadline := time.Now().Add(5 * time.Second)
	for heartbeatGrows(file) {
		if time.Now().After(deadline) {
			t.Fatal("grandchild survived the timeout")
		}
	}
}

func TestRunCommandCancellationKillsProcessTree(t *testing.T) {
	file := filepath.Join(t.TempDir(), "heartbeat")
	t.Cleanup(func() { killHelperByPIDFile(file) })
	name, args := shellHelper(t, "background-hang", file)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			if info, err := os.Stat(file); err == nil && info.Size() > 0 {
				cancel()
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	_, err := runCommand(ctx, t.TempDir(), name, args, 60)
	if err == nil || !strings.Contains(err.Error(), "canceled") {
		t.Fatalf("err = %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for heartbeatGrows(file) {
		if time.Now().After(deadline) {
			t.Fatal("grandchild survived cancellation")
		}
	}
}

func TestRunCommandBoundsOutputMemory(t *testing.T) {
	name, args := shellHelper(t, "flood")
	output, err := runCommand(context.Background(), t.TempDir(), name, args, 60)
	if err != nil {
		t.Fatal(err)
	}
	if len(output) > maxToolOutputBytes+256 {
		t.Fatalf("output is %d bytes, limit %d", len(output), maxToolOutputBytes)
	}
	if !strings.Contains(output, "output truncated") {
		t.Fatalf("missing truncation marker: %q", output[len(output)-200:])
	}
}

func TestCappedBufferKeepsPrefix(t *testing.T) {
	b := &cappedBuffer{limit: 5}
	for _, part := range []string{"ab", "cdef", "gh"} {
		if n, err := b.Write([]byte(part)); n != len(part) || err != nil {
			t.Fatalf("Write(%q) = %d, %v", part, n, err)
		}
	}
	if got := b.String(); !strings.HasPrefix(got, "abcde\n") || !strings.Contains(got, "3 more bytes") {
		t.Fatalf("String() = %q", got)
	}
}
