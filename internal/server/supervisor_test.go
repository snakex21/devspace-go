package server

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func exitedTunnel() *tunnelProcess {
	done := make(chan struct{})
	close(done)
	return &tunnelProcess{url: "https://fixture.trycloudflare.com", done: done, cancel: func() {}, err: errors.New("exit status 1")}
}

func TestBackoffDelayGrowsWithJitterAndCap(t *testing.T) {
	policy := backoffPolicy{base: time.Second, max: 8 * time.Second}
	low := func(int64) int64 { return 0 }
	high := func(n int64) int64 { return n - 1 }
	for attempt, want := range []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 8 * time.Second, 8 * time.Second} {
		if got := policy.delay(attempt, high); got != want {
			t.Errorf("attempt %d max delay = %v, want %v", attempt, got, want)
		}
		if got := policy.delay(attempt, low); got != want/2 {
			t.Errorf("attempt %d min delay = %v, want %v", attempt, got, want/2)
		}
	}
	// Huge attempt counts must not overflow into negative or zero delays.
	if got := policy.delay(10_000, high); got != 8*time.Second {
		t.Fatalf("delay after many attempts = %v", got)
	}
}

func TestSuperviseTunnelRestartsExitedTunnelWithBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var mu sync.Mutex
	var starts []time.Time
	start := func(context.Context, string) (*tunnelProcess, error) {
		mu.Lock()
		starts = append(starts, time.Now())
		n := len(starts)
		mu.Unlock()
		if n == 4 {
			cancel()
		}
		return exitedTunnel(), nil
	}
	done := make(chan struct{})
	go func() {
		superviseTunnel(ctx, "http://127.0.0.1:1", start, backoffPolicy{base: 40 * time.Millisecond, max: 160 * time.Millisecond, resetAfter: time.Hour})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("supervisor did not stop after cancellation")
	}
	mu.Lock()
	defer mu.Unlock()
	if len(starts) != 4 {
		t.Fatalf("tunnel started %d times, want 4", len(starts))
	}
	// Minimum (jittered) gaps are 20ms, 40ms, 80ms: no tight loop.
	for i, minimum := range []time.Duration{20 * time.Millisecond, 40 * time.Millisecond, 80 * time.Millisecond} {
		if gap := starts[i+1].Sub(starts[i]); gap < minimum {
			t.Errorf("restart %d after %v, want >= %v", i+1, gap, minimum)
		}
	}
}

func TestSuperviseTunnelRetriesFailedStartup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var calls atomic.Int32
	start := func(context.Context, string) (*tunnelProcess, error) {
		if calls.Add(1) == 3 {
			cancel()
		}
		return nil, errors.New("temporary DNS failure")
	}
	superviseTunnel(ctx, "http://127.0.0.1:1", start, backoffPolicy{base: time.Millisecond, max: 2 * time.Millisecond})
	if calls.Load() != 3 {
		t.Fatalf("startup attempted %d times, want 3", calls.Load())
	}
}

func TestSuperviseTunnelStopsWhenNoProviderInstalled(t *testing.T) {
	var calls atomic.Int32
	start := func(context.Context, string) (*tunnelProcess, error) {
		calls.Add(1)
		return nil, errNoTunnelProvider
	}
	finished := make(chan struct{})
	go func() {
		superviseTunnel(context.Background(), "http://127.0.0.1:1", start, backoffPolicy{base: time.Millisecond, max: time.Millisecond})
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("supervisor kept retrying without any tunnel program")
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d", calls.Load())
	}
}

func TestSuperviseTunnelCancelDuringBackoffReturnsPromptly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{}, 1)
	start := func(context.Context, string) (*tunnelProcess, error) {
		started <- struct{}{}
		return nil, errors.New("fail")
	}
	finished := make(chan struct{})
	go func() {
		superviseTunnel(ctx, "http://127.0.0.1:1", start, backoffPolicy{base: time.Hour, max: time.Hour})
		close(finished)
	}()
	<-started
	cancel()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("supervisor slept through cancellation")
	}
}

func TestSuperviseTunnelStopsRunningTunnelOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	done := make(chan struct{})
	tunnel := &tunnelProcess{url: "https://fixture.trycloudflare.com", done: done, cancel: func() { close(stopped); close(done) }}
	finished := make(chan struct{})
	go func() {
		superviseTunnel(ctx, "http://127.0.0.1:1", func(context.Context, string) (*tunnelProcess, error) { return tunnel, nil }, defaultTunnelBackoff)
		close(finished)
	}()
	cancel()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("supervisor did not return")
	}
	select {
	case <-stopped:
	default:
		t.Fatal("running tunnel was not stopped")
	}
}
