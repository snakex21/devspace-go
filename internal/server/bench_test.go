package server

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/store"
	"github.com/snakex21/devspace-go/internal/workspace"
)

func BenchmarkCreateMcpServer(b *testing.B) {
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{b.TempDir()}
	s := &Server{cfg: cfg, registry: workspace.NewRegistry(cfg, nil)}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.createMcpServer()
	}
}

// BenchmarkMCPReadToolCall measures one complete stateless /mcp tool call
// (HTTP + JSON-RPC + workspace lookup + read), the server's hot path.
func BenchmarkMCPReadToolCall(b *testing.B) {
	root := b.TempDir()
	if err := os.WriteFile(filepath.Join(root, "fixture.txt"), []byte("hello"), 0o600); err != nil {
		b.Fatal(err)
	}
	st, err := store.New(b.TempDir())
	if err != nil {
		b.Fatal(err)
	}
	defer st.Close()
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}
	cfg.Logging.Requests = false
	s := &Server{cfg: cfg, registry: workspace.NewRegistry(cfg, st), store: st}
	httpServer := httptest.NewServer(s.streamableMCPHandler())
	defer httpServer.Close()
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "bench", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer session.Close()
	params := &mcp.CallToolParams{Name: "read", Arguments: map[string]any{"workspaceId": "default", "path": "fixture.txt"}}
	if result, err := session.CallTool(ctx, params); err != nil || result.IsError {
		b.Fatalf("%v %#v", err, result)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := session.CallTool(ctx, params); err != nil {
			b.Fatal(err)
		}
	}
}
