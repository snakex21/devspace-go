package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/workspace"
)

var tunnelHeaders = map[string]string{"Cf-Connecting-Ip": "203.0.113.9", "Cf-Ray": "1"}

func fullAccessServer(t *testing.T, remote *bool) (s *Server, endpoint, root, outside string) {
	t.Helper()
	parent := t.TempDir()
	root = filepath.Join(parent, "ws")
	outside = filepath.Join(parent, "elsewhere")
	for _, dir := range []string{root, outside} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(outside, "notes.txt"), []byte("OUTSIDE-MARKER"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}
	cfg.StateDir = filepath.Join(parent, "state")
	cfg.ConfigDir = filepath.Join(parent, "config")
	cfg.OutsideAccess.Remote = remote
	policy, problems := workspace.NewAccessPolicy(cfg)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	registry := workspace.NewRegistry(cfg, nil)
	registry.SetAccessPolicy(policy)
	s = &Server{cfg: cfg, registry: registry, policy: policy}
	httpServer := httptest.NewServer(s.streamableMCPHandler())
	t.Cleanup(httpServer.Close)
	return s, httpServer.URL, root, outside
}

func mcpSession(t *testing.T, endpoint string, header map[string]string) *mcp.ClientSession {
	t.Helper()
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		clone := req.Clone(req.Context())
		for name, value := range header {
			clone.Header.Set(name, value)
		}
		return http.DefaultTransport.RoundTrip(clone)
	})}
	client := mcp.NewClient(&mcp.Implementation{Name: "full-access-test", Version: "1.0.0"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: endpoint, HTTPClient: httpClient}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

func callTool(t *testing.T, session *mcp.ClientSession, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestMCPFullAccessByDefaultLocallyAndThroughTunnel(t *testing.T) {
	s, endpoint, _, outside := fullAccessServer(t, nil)
	for _, header := range []map[string]string{nil, tunnelHeaders} {
		session := mcpSession(t, endpoint, header)
		if r := callTool(t, session, "open_default_workspace", map[string]any{}); r.IsError {
			t.Fatalf("open_default_workspace: %s", textOf(r))
		}
		read := callTool(t, session, "read", map[string]any{"workspaceId": "default", "path": filepath.Join(outside, "notes.txt")})
		if read.IsError || !strings.Contains(textOf(read), "OUTSIDE-MARKER") {
			t.Fatalf("read outside (tunnel=%v): %s", header != nil, textOf(read))
		}
		write := callTool(t, session, "write", map[string]any{"workspaceId": "default", "path": filepath.Join(outside, "w.txt"), "content": "x"})
		if write.IsError {
			t.Fatalf("write outside (tunnel=%v): %s", header != nil, textOf(write))
		}
		// open_workspace accepts a folder outside allowedRoots.
		opened := callTool(t, session, "open_workspace", map[string]any{"path": outside})
		if opened.IsError {
			t.Fatalf("open_workspace outside allowedRoots (tunnel=%v): %s", header != nil, textOf(opened))
		}
	}
	for _, origin := range []workspace.Origin{workspace.OriginLocal, workspace.OriginRemote} {
		if !strings.Contains(s.serverInstructions(origin), "Full disk access is enabled") {
			t.Fatalf("%s instructions do not mention full access", origin)
		}
	}
}

func TestMCPFullAccessRemoteDisabled(t *testing.T) {
	off := false
	s, endpoint, _, outside := fullAccessServer(t, &off)

	local := mcpSession(t, endpoint, nil)
	callTool(t, local, "open_default_workspace", map[string]any{})
	if r := callTool(t, local, "read", map[string]any{"workspaceId": "default", "path": filepath.Join(outside, "notes.txt")}); r.IsError {
		t.Fatalf("local read outside: %s", textOf(r))
	}

	remote := mcpSession(t, endpoint, tunnelHeaders)
	opened := callTool(t, remote, "open_default_workspace", map[string]any{})
	if opened.IsError {
		t.Fatalf("remote open_default_workspace: %s", textOf(opened))
	}
	read := callTool(t, remote, "read", map[string]any{"workspaceId": "default", "path": filepath.Join(outside, "notes.txt")})
	if !read.IsError || strings.Contains(textOf(read), "OUTSIDE-MARKER") {
		t.Fatalf("remote read outside with remote:false: %s", textOf(read))
	}
	if r := callTool(t, remote, "open_workspace", map[string]any{"path": outside}); !r.IsError {
		t.Fatalf("remote open_workspace outside allowedRoots with remote:false: %s", textOf(r))
	}
	if strings.Contains(s.serverInstructions(workspace.OriginRemote), "Full disk access") {
		t.Fatal("remote instructions announce full access although it is disabled")
	}
}
