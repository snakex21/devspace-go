package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/workspace"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestStreamableMCPStaysUsableAcrossRepeatedRequests(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{t.TempDir()}

	server := &Server{
		cfg:      cfg,
		registry: workspace.NewRegistry(cfg, nil),
	}
	httpServer := httptest.NewServer(server.streamableMCPHandler())
	defer httpServer.Close()

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "regression-test", Version: "1.0.0"}, nil)
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			clone := req.Clone(req.Context())
			clone.Header.Del("Mcp-Session-Id")
			return http.DefaultTransport.RoundTrip(clone)
		}),
	}
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:   httpServer.URL,
		HTTPClient: httpClient,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	for request := 1; request <= 25; request++ {
		result, err := session.ListTools(ctx, nil)
		if err != nil {
			t.Fatalf("request %d failed: %v", request, err)
		}
		if len(result.Tools) < 2 {
			t.Fatalf("request %d returned only %d tools", request, len(result.Tools))
		}
	}

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "open_default_workspace",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("open_default_workspace returned an MCP error: %#v", result.Content)
	}
}
