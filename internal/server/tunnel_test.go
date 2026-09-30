package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Re-execute only this helper as a local fake tunnel. It never opens a network
// connection and runs on Windows as well as Unix.
func TestTunnelHelperProcess(t *testing.T) {
	if os.Getenv("DEVSPACE_TUNNEL_HELPER") != "1" {
		return
	}
	args := os.Args
	for len(args) > 0 && args[0] != "--" {
		args = args[1:]
	}
	mode, marker := args[1], args[2]
	switch mode {
	case "exit":
		os.Exit(7)
	case "hang":
		time.Sleep(time.Minute)
		os.Exit(0)
	case "oversized":
		fmt.Fprintln(os.Stderr, strings.Repeat("x", maxTunnelLogLine*2))
	}
	fmt.Fprintln(os.Stdout, "https://fixture.trycloudflare.com")
	fmt.Fprintln(os.Stderr, "https://fixture.trycloudflare.com")
	if mode == "url-exit" {
		os.Exit(7)
	}
	if mode == "flood" {
		for i := 0; i < 10000; i++ {
			fmt.Fprintln(os.Stdout, strings.Repeat("x", 128))
			fmt.Fprintln(os.Stderr, strings.Repeat("y", 128))
		}
	}
	if err := os.WriteFile(marker, []byte("drained"), 0600); err != nil {
		os.Exit(8)
	}
	time.Sleep(time.Minute)
	os.Exit(0)
}

func helperCommand(mode, marker string, command **exec.Cmd) func(context.Context) *exec.Cmd {
	return func(ctx context.Context) *exec.Cmd {
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestTunnelHelperProcess$", "--", mode, marker)
		cmd.Env = append(os.Environ(), "DEVSPACE_TUNNEL_HELPER=1")
		*command = cmd
		return cmd
	}
}

func waitForMarker(ctx context.Context, marker string) error {
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := os.Stat(marker); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func TestTunnelDrainsBothStreamsAfterURLAndReapsOnStop(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "drained")
	var cmd *exec.Cmd
	process, err := startTunnelProcess(context.Background(), helperCommand("flood", marker, &cmd), cloudflaredURLPattern, 10*time.Second, func(ctx context.Context, _ string) error { return waitForMarker(ctx, marker) }, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(process.stop)
	if process.url != "https://fixture.trycloudflare.com" {
		t.Fatalf("URL = %q", process.url)
	}
	process.stop()
	// stop is safe to call twice, and does not return before Wait finishes.
	process.stop()
	if cmd.ProcessState == nil {
		t.Fatal("child was not reaped")
	}
}

func TestTunnelOversizedLogDoesNotStopURLDetection(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "drained")
	var cmd *exec.Cmd
	process, err := startTunnelProcess(context.Background(), helperCommand("oversized", marker, &cmd), cloudflaredURLPattern, 10*time.Second, func(ctx context.Context, _ string) error { return waitForMarker(ctx, marker) }, nil)
	if err != nil {
		t.Fatal(err)
	}
	process.stop()
}

func TestTunnelStartupFailureAlwaysReapsProcess(t *testing.T) {
	for _, mode := range []string{"exit", "url-exit", "hang"} {
		t.Run(mode, func(t *testing.T) {
			var cmd *exec.Cmd
			timeout := 5 * time.Second
			if mode == "hang" {
				timeout = 150 * time.Millisecond
			}
			started := time.Now()
			process, err := startTunnelProcess(context.Background(), helperCommand(mode, filepath.Join(t.TempDir(), "marker"), &cmd), cloudflaredURLPattern, timeout, func(ctx context.Context, _ string) error { <-ctx.Done(); return ctx.Err() }, nil)
			if err == nil || process != nil {
				t.Fatalf("got process=%#v error=%v", process, err)
			}
			if cmd.ProcessState == nil {
				t.Fatal("failed child was not reaped")
			}
			if time.Since(started) > 3*time.Second {
				t.Fatal("waited for URL timeout after process exited")
			}
		})
	}
}

func TestTunnelCancellationBeforeURLAndDuringReadiness(t *testing.T) {
	for _, mode := range []string{"hang", "ready"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			marker := filepath.Join(t.TempDir(), "marker")
			var cmd *exec.Cmd
			finished := make(chan error, 1)
			startedProbe := make(chan struct{})
			go func() {
				_, err := startTunnelProcess(ctx, helperCommand(mode, marker, &cmd), cloudflaredURLPattern, time.Minute, func(ctx context.Context, _ string) error {
					close(startedProbe)
					<-ctx.Done()
					return ctx.Err()
				}, nil)
				finished <- err
			}()
			if mode == "ready" {
				select {
				case <-startedProbe:
				case <-time.After(5 * time.Second):
					t.Fatal("readiness probe never started")
				}
			}
			cancel()
			select {
			case err := <-finished:
				if err == nil {
					t.Fatal("canceled startup succeeded")
				}
				if cmd.Process != nil && cmd.ProcessState == nil {
					t.Fatal("canceled child was not reaped")
				}
			case <-time.After(3 * time.Second):
				t.Fatal("cancel did not interrupt startup")
			}
		})
	}
}

func TestTunnelProvidersStopRetryingOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	provider := func(context.Context, string) (*tunnelProcess, error) { calls++; cancel(); return nil, context.Canceled }
	s := &Server{}
	if p := s.startTunnelWithProviders(ctx, "http://127.0.0.1:1", provider, provider); p != nil {
		t.Fatal("canceled provider returned a process")
	}
	if calls != 1 {
		t.Fatalf("called providers %d times after cancellation", calls)
	}
}

func TestTunnelHealthWaitsThroughDNSAndGatewayFailures(t *testing.T) {
	var requests atomic.Int32
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("probe path = %q", r.URL.Path)
		}
		switch requests.Add(1) {
		case 1:
			http.Error(w, "not ready", http.StatusBadGateway)
		case 2:
			fmt.Fprint(w, `{"ok":true,"name":"another-service"}`)
		default:
			fmt.Fprint(w, `{"ok":true,"name":"devspace-go"}`)
		}
	}))
	defer origin.Close()
	client := origin.Client()
	transport := client.Transport
	var failures atomic.Int32
	client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if failures.Add(1) == 1 {
			return nil, errors.New("temporary DNS error")
		}
		return transport.RoundTrip(req)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := waitForTunnelHealthWithClient(ctx, origin.URL, client, time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if requests.Load() != 3 {
		t.Fatalf("health requests = %d", requests.Load())
	}
}

func TestTunnelHealthNeverAcceptsWrongOrUnavailableOrigin(t *testing.T) {
	for _, response := range []string{`{"ok":false,"name":"devspace-go"}`, `{"ok":true,"name":"other"}`, `<html>Login</html>`} {
		t.Run(response, func(t *testing.T) {
			origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, response) }))
			defer origin.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
			defer cancel()
			err := waitForTunnelHealthWithClient(ctx, origin.URL, origin.Client(), time.Millisecond)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("got %v", err)
			}
		})
	}
}

func TestTunnelHealthCancellationInterruptsInFlightRequest(t *testing.T) {
	entered := make(chan struct{})
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(entered); <-r.Context().Done() }))
	defer origin.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- waitForTunnelHealthWithClient(ctx, origin.URL, origin.Client(), time.Millisecond) }()
	<-entered
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("probe did not cancel")
	}
}

func TestTunnelOrigin(t *testing.T) {
	for address, want := range map[string]string{"0.0.0.0:7676": "http://127.0.0.1:7676", "[::]:7676": "http://[::1]:7676", "127.0.0.1:1234": "http://127.0.0.1:1234", "[::1]:1234": "http://[::1]:1234"} {
		if got := tunnelOrigin(address); got != want {
			t.Errorf("%s: got %s, want %s", address, got, want)
		}
	}
}

func TestTunnelOutputHandlesSplitLinesAndOversizedLines(t *testing.T) {
	var lines []string
	writer := &tunnelOutput{onLine: func(line string) { lines = append(lines, line) }}
	for _, part := range []string{"https://fixture.", "trycloudflare.com\r\n", strings.Repeat("x", maxTunnelLogLine+1), "\nnext\n"} {
		if n, err := writer.Write([]byte(part)); err != nil || n != len(part) {
			t.Fatalf("Write=%d,%v", n, err)
		}
	}
	if len(lines) != 2 || lines[0] != "https://fixture.trycloudflare.com" || lines[1] != "next" {
		t.Fatalf("lines=%#v", lines)
	}
}

func TestFindCloudflaredExecutableUsesPlatformBinary(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	tools := filepath.Join(dir, "tools")
	if err := os.Mkdir(tools, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"cloudflared", "cloudflared.exe"} {
		if err := os.WriteFile(filepath.Join(tools, name), []byte("fixture"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	name := "cloudflared"
	if runtime.GOOS == "windows" {
		name = "cloudflared.exe"
	}
	if got := findCloudflaredExecutable(); got != filepath.Join(tools, name) {
		t.Fatalf("found %q, want platform binary %q", got, name)
	}
}
