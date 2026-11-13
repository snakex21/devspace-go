package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/workspace"
)

type panicInput struct {
	Value string `json:"value"`
}

type panicOutput struct {
	Result string `json:"result"`
}

func TestToolPanicIsReturnedAsErrorAndServerSurvives(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "panic-test", Version: "1"}, nil)
	addTool(server, &mcp.Tool{Name: "explode"}, func(ctx context.Context, req *mcp.CallToolRequest, input panicInput) (*mcp.CallToolResult, panicOutput, error) {
		if input.Value == "boom" {
			var m map[string]int
			m["nil map write"] = 1 // runtime panic on the SDK's handler goroutine
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "fine"}}}, panicOutput{Result: "fine"}, nil
	})
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	ctx := context.Background()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "explode", Arguments: map[string]any{"value": "boom"}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError || !strings.Contains(textOf(result), "internal error in tool explode") {
		t.Fatalf("panic result = %#v", result)
	}
	result, err = session.CallTool(ctx, &mcp.CallToolParams{Name: "explode", Arguments: map[string]any{"value": "ok"}})
	if err != nil || result.IsError {
		t.Fatalf("server unusable after panic: %v %#v", err, result)
	}
}

func textOf(result *mcp.CallToolResult) string {
	var parts []string
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			parts = append(parts, text.Text)
		}
	}
	return strings.Join(parts, "\n")
}

func TestRecoverMiddlewareReturns500(t *testing.T) {
	handler := recoverMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/mcp", nil))
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", recorder.Code)
	}
}

func TestRequestBodyLimit(t *testing.T) {
	var readErr error
	handler := limitRequestBody(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, readErr = io.Copy(io.Discard, r.Body)
	}))

	declared := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader("x"))
	declared.ContentLength = maxRequestBodyBytes + 1
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, declared)
	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("declared oversize status = %d", recorder.Code)
	}

	// A chunked body has no Content-Length; reading must still stop.
	chunked := httptest.NewRequest(http.MethodPost, "/mcp", io.LimitReader(zeroReader{}, maxRequestBodyBytes+1024))
	chunked.ContentLength = -1
	handler.ServeHTTP(httptest.NewRecorder(), chunked)
	if readErr == nil || !strings.Contains(readErr.Error(), "too large") {
		t.Fatalf("chunked oversize read error = %v", readErr)
	}
}

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	clear(p)
	return len(p), nil
}

func TestHTTPServerHasTimeouts(t *testing.T) {
	server := newHTTPServer(http.NotFoundHandler(), context.Background())
	if server.ReadHeaderTimeout <= 0 || server.ReadTimeout <= 0 || server.IdleTimeout <= 0 || server.MaxHeaderBytes <= 0 {
		t.Fatalf("missing timeouts: %+v", server)
	}
}

func TestSlowHeaderConnectionIsClosed(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := newHTTPServer(http.NotFoundHandler(), context.Background())
	server.ReadHeaderTimeout = 200 * time.Millisecond
	go server.Serve(listener)
	defer server.Close()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.Write([]byte("GET / HTTP/1.1\r\nHost: x\r\n")) // never finish the headers
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	started := time.Now()
	io.Copy(io.Discard, conn)
	if elapsed := time.Since(started); elapsed > 3*time.Second {
		t.Fatalf("slow-header connection held open for %v", elapsed)
	}
}

func TestListenErrorNamesRunningInstance(t *testing.T) {
	first := lifecycleServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	originCh := make(chan string, 1)
	done := make(chan error, 1)
	go func() {
		done <- first.start(ctx, func(_ context.Context, origin string) (*tunnelProcess, error) {
			originCh <- origin
			return nil, errNoTunnelProvider
		})
	}()
	origin := <-originCh
	_, port, _ := net.SplitHostPort(strings.TrimPrefix(origin, "http://"))

	second := lifecycleServer(t)
	second.cfg.Port, _ = strconv.Atoi(port)
	err := second.start(context.Background(), func(context.Context, string) (*tunnelProcess, error) {
		t.Error("tunnel started for a server that could not bind")
		return nil, errNoTunnelProvider
	})
	if err == nil || !strings.Contains(err.Error(), "another Dev Space Go instance") {
		t.Fatalf("err = %v", err)
	}
	cancel()
	<-done
}

func TestListenErrorNamesForeignProgram(t *testing.T) {
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "hello") }))
	defer foreign.Close()
	_, port, _ := net.SplitHostPort(foreign.Listener.Addr().String())
	server := lifecycleServer(t)
	server.cfg.Port, _ = strconv.Atoi(port)
	err := server.start(context.Background(), func(context.Context, string) (*tunnelProcess, error) { return nil, errNoTunnelProvider })
	if err == nil || !strings.Contains(err.Error(), "another program") || !strings.Contains(err.Error(), "listen") {
		t.Fatalf("err = %v", err)
	}
}

func TestInvalidPortIsRejected(t *testing.T) {
	server := lifecycleServer(t)
	server.cfg.Port = 70000
	err := server.start(context.Background(), func(context.Context, string) (*tunnelProcess, error) { return nil, errNoTunnelProvider })
	if err == nil || !strings.Contains(err.Error(), "invalid port") {
		t.Fatalf("err = %v", err)
	}
}

func TestShutdownDoesNotWaitForOpenSSEStream(t *testing.T) {
	server := lifecycleServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	originCh := make(chan string, 1)
	done := make(chan error, 1)
	go func() {
		done <- server.start(ctx, func(_ context.Context, origin string) (*tunnelProcess, error) {
			originCh <- origin
			return nil, errNoTunnelProvider
		})
	}()
	origin := <-originCh
	response, err := http.Get(origin + "/sse")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	// Wait for the endpoint event so the stream is definitely established.
	if _, err := bufio.NewReader(response.Body).ReadString('\n'); err != nil {
		t.Fatal(err)
	}

	started := time.Now()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(shutdownTimeout + 5*time.Second):
		t.Fatal("server did not stop")
	}
	if elapsed := time.Since(started); elapsed > 3*time.Second {
		t.Fatalf("shutdown waited %v for an idle SSE stream", elapsed)
	}
}

func TestMCPToolsRejectTraversalEndToEnd(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "ws")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(parent, "secret.txt"), []byte("TOP-SECRET-MARKER"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}
	server := &Server{cfg: cfg, registry: workspace.NewRegistry(cfg, nil)}
	httpServer := httptest.NewServer(server.streamableMCPHandler())
	defer httpServer.Close()

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "traversal-test", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	for _, call := range []mcp.CallToolParams{
		{Name: "ls", Arguments: map[string]any{"workspaceId": "default", "path": ".."}},
		{Name: "grep", Arguments: map[string]any{"workspaceId": "default", "pattern": "TOP-SECRET", "path": ".."}},
		{Name: "glob", Arguments: map[string]any{"workspaceId": "default", "pattern": "*.txt", "path": ".."}},
		{Name: "read", Arguments: map[string]any{"workspaceId": "default", "path": "../secret.txt"}},
		{Name: "move", Arguments: map[string]any{"workspaceId": "default", "sourcePath": "a", "targetPath": ".", "overwrite": true}},
	} {
		result, err := session.CallTool(ctx, &call)
		if err != nil {
			t.Fatalf("%s: %v", call.Name, err)
		}
		if !result.IsError || strings.Contains(textOf(result), "TOP-SECRET-MARKER") || strings.Contains(textOf(result), "secret.txt\n") {
			t.Fatalf("%s escaped the workspace: %s", call.Name, textOf(result))
		}
	}
}
