package server

import (
	"bufio"
	"context"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestServerLifecycleLeavesNoGoroutines runs complete start/serve/stop
// cycles, including MCP and SSE traffic, and checks that every goroutine the
// server started (listeners, tunnel supervisor, SSE streams, MCP sessions)
// has exited afterwards.
func TestServerLifecycleLeavesNoGoroutines(t *testing.T) {
	http.DefaultClient.CloseIdleConnections()
	baseline := runtime.NumGoroutine()

	for cycle := 0; cycle < 3; cycle++ {
		server := lifecycleServer(t)
		ctx, cancel := context.WithCancel(context.Background())
		originCh := make(chan string, 1)
		done := make(chan error, 1)
		go func() {
			done <- server.start(ctx, func(ctx context.Context, origin string) (*tunnelProcess, error) {
				originCh <- origin
				<-ctx.Done()
				return nil, ctx.Err()
			})
		}()
		origin := <-originCh

		transport := &http.Transport{}
		client := mcp.NewClient(&mcp.Implementation{Name: "leak-test", Version: "1"}, nil)
		session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: origin + "/mcp", HTTPClient: &http.Client{Transport: transport}}, nil)
		if err != nil {
			t.Fatal(err)
		}
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "open_default_workspace", Arguments: map[string]any{}})
		if err != nil || result.IsError {
			t.Fatalf("call failed: %v %#v", err, result)
		}
		session.Close()

		sse, err := (&http.Client{Transport: transport}).Get(origin + "/sse")
		if err != nil {
			t.Fatal(err)
		}
		bufio.NewReader(sse.Body).ReadString('\n')

		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(15 * time.Second):
			t.Fatal("server did not stop")
		}
		sse.Body.Close()
		transport.CloseIdleConnections()
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		current := runtime.NumGoroutine()
		if current <= baseline {
			return
		}
		if time.Now().After(deadline) {
			buf := make([]byte, 1<<20)
			n := runtime.Stack(buf, true)
			t.Fatalf("goroutines: baseline %d, now %d\n%s", baseline, current, buf[:n])
		}
		time.Sleep(20 * time.Millisecond)
	}
}
