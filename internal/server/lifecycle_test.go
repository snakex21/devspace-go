package server

import (
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/workspace"
)

func lifecycleServer(t *testing.T) *Server {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = 0
	cfg.AllowedRoots = []string{t.TempDir()}
	return &Server{cfg: cfg, registry: workspace.NewRegistry(cfg, nil)}
}

func TestServerServesBeforeTunnelReadyAndCancelsStartup(t *testing.T) {
	server := lifecycleServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	originReady := make(chan string, 1)
	stopped := make(chan struct{})
	finished := make(chan error, 1)
	go func() {
		finished <- server.start(ctx, func(ctx context.Context, origin string) *tunnelProcess {
			originReady <- origin
			<-ctx.Done()
			close(stopped)
			return nil
		})
	}()
	var origin string
	select {
	case origin = <-originReady:
	case <-time.After(3 * time.Second):
		t.Fatal("tunnel startup never began")
	}
	client := &http.Client{Timeout: time.Second}
	response, err := client.Get(origin + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil || response.StatusCode != http.StatusOK || !strings.Contains(string(body), `"name":"devspace-go"`) {
		t.Fatalf("health = %d %q %v", response.StatusCode, body, err)
	}
	if err := server.start(ctx, func(context.Context, string) *tunnelProcess { t.Error("second start launched a tunnel"); return nil }); err == nil {
		t.Fatal("same Server started twice")
	}
	cancel()
	select {
	case err := <-finished:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("server did not stop during tunnel startup")
	}
	select {
	case <-stopped:
	default:
		t.Fatal("server returned before tunnel startup was stopped")
	}
	if _, err := client.Get(origin + "/healthz"); err == nil {
		t.Fatal("listener remained open after shutdown")
	}
}

func TestServerBindFailureNeverStartsTunnel(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	server := lifecycleServer(t)
	server.cfg.Port = listener.Addr().(*net.TCPAddr).Port
	err = server.start(context.Background(), func(context.Context, string) *tunnelProcess {
		t.Error("tunnel launched despite occupied port")
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "listen") {
		t.Fatalf("got %v", err)
	}
}

func TestServerRepeatedNewInstancesReleasePort(t *testing.T) {
	port := 0
	for run := 0; run < 5; run++ {
		server := lifecycleServer(t)
		server.cfg.Port = port
		ctx, cancel := context.WithCancel(context.Background())
		entered := make(chan string, 1)
		done := make(chan error, 1)
		go func() {
			done <- server.start(ctx, func(_ context.Context, origin string) *tunnelProcess { entered <- origin; return nil })
		}()
		var origin string
		select {
		case origin = <-entered:
		case err := <-done:
			cancel()
			t.Fatalf("run %d failed: %v", run, err)
		case <-time.After(3 * time.Second):
			cancel()
			t.Fatal("start timed out")
		}
		_, portText, err := net.SplitHostPort(strings.TrimPrefix(origin, "http://"))
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		port, err = strconv.Atoi(portText)
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("shutdown timed out")
		}
	}
}

func TestServerCanceledContextDoesNotStartTunnel(t *testing.T) {
	server := lifecycleServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := server.start(ctx, func(context.Context, string) *tunnelProcess {
		t.Error("tunnel launched with canceled context")
		return nil
	}); err != context.Canceled {
		t.Fatalf("got %v", err)
	}
}
