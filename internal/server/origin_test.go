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

func TestRequestOriginIsLocalOnlyWithoutAnyTunnelSignal(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		host       string
		header     map[string]string
		listenHost string
		want       workspace.Origin
	}{
		{name: "loopback ipv4", remoteAddr: "127.0.0.1:50000", host: "127.0.0.1:7676", want: workspace.OriginLocal},
		{name: "localhost", remoteAddr: "127.0.0.1:50000", host: "localhost:7676", want: workspace.OriginLocal},
		{name: "loopback ipv6", remoteAddr: "[::1]:50000", host: "[::1]:7676", want: workspace.OriginLocal},
		{name: "cloudflare host", remoteAddr: "127.0.0.1:50000", host: "abc.trycloudflare.com", want: workspace.OriginRemote},
		{name: "pinggy host", remoteAddr: "127.0.0.1:50000", host: "abc.a.free.pinggy.link", want: workspace.OriginRemote},
		{name: "cloudflare headers with rewritten host", remoteAddr: "127.0.0.1:50000", host: "127.0.0.1:7676", header: map[string]string{"Cf-Connecting-Ip": "203.0.113.9"}, want: workspace.OriginRemote},
		{name: "cf-ray", remoteAddr: "127.0.0.1:50000", host: "127.0.0.1:7676", header: map[string]string{"Cf-Ray": "x"}, want: workspace.OriginRemote},
		{name: "x-forwarded-for", remoteAddr: "127.0.0.1:50000", host: "127.0.0.1:7676", header: map[string]string{"X-Forwarded-For": "203.0.113.9"}, want: workspace.OriginRemote},
		{name: "forwarded", remoteAddr: "127.0.0.1:50000", host: "127.0.0.1:7676", header: map[string]string{"Forwarded": "for=203.0.113.9"}, want: workspace.OriginRemote},
		{name: "lan peer", remoteAddr: "192.168.1.20:50000", host: "192.168.1.10:7676", listenHost: "0.0.0.0", want: workspace.OriginRemote},
		{name: "lan peer claiming localhost", remoteAddr: "192.168.1.20:50000", host: "127.0.0.1:7676", want: workspace.OriginRemote},
		{name: "missing host", remoteAddr: "127.0.0.1:50000", host: "", want: workspace.OriginRemote},
		{name: "garbage remote addr", remoteAddr: "nonsense", host: "127.0.0.1:7676", want: workspace.OriginRemote},
		{name: "dns rebinding name", remoteAddr: "127.0.0.1:50000", host: "evil.example:7676", want: workspace.OriginRemote},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
			r.RemoteAddr = test.remoteAddr
			r.Host = test.host
			for name, value := range test.header {
				r.Header.Set(name, value)
			}
			listenHost := test.listenHost
			if listenHost == "" {
				listenHost = "127.0.0.1"
			}
			if got := requestOrigin(r, listenHost); got != test.want {
				t.Fatalf("requestOrigin = %s, want %s", got, test.want)
			}
		})
	}
	if requestOrigin(nil, "127.0.0.1") != workspace.OriginRemote {
		t.Fatal("nil request must be remote")
	}
}

// TestMCPExtraRootsLocalVersusTunnel drives the real /mcp handler: the same
// read of an extra root works for a local client and is refused when the
// request carries the headers a Cloudflare tunnel adds only when remote
// access to extra roots is explicitly disabled.
func TestMCPExtraRootsLocalVersusTunnel(t *testing.T) {
	off, on := false, true
	// nil is the default: remote access to extra roots is on.
	for _, setting := range []*bool{nil, &off, &on} {
		remoteEnabled := setting == nil || *setting
		parent := t.TempDir()
		root := filepath.Join(parent, "ws")
		shared := filepath.Join(parent, "shared")
		for _, dir := range []string{root, shared} {
			if err := os.Mkdir(dir, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(shared, "lib.go"), []byte("SHARED-MARKER"), 0o600); err != nil {
			t.Fatal(err)
		}
		cfg := config.DefaultConfig()
		cfg.AllowedRoots = []string{root}
		cfg.StateDir = filepath.Join(parent, "state")
		cfg.ConfigDir = filepath.Join(parent, "config")
		cfg.OutsideAccess = config.OutsideAccessConfig{
			Mode:   config.AccessModeRoots,
			Roots:  []config.ExtraRootConfig{{Path: shared, Access: "read"}},
			Remote: setting,
		}
		policy, problems := workspace.NewAccessPolicy(cfg)
		if len(problems) > 0 {
			t.Fatal(problems)
		}
		s := &Server{cfg: cfg, registry: workspace.NewRegistry(cfg, nil), policy: policy}
		httpServer := httptest.NewServer(s.streamableMCPHandler())

		local := callRead(t, httpServer.URL, nil, filepath.Join(shared, "lib.go"))
		if local.IsError || !strings.Contains(textOf(local), "SHARED-MARKER") {
			t.Fatalf("local read of extra root: %s", textOf(local))
		}
		tunnel := callRead(t, httpServer.URL, map[string]string{"Cf-Connecting-Ip": "203.0.113.9", "Cf-Ray": "1"}, filepath.Join(shared, "lib.go"))
		if remoteEnabled {
			if tunnel.IsError {
				t.Fatalf("tunnel read with remote enabled: %s", textOf(tunnel))
			}
		} else if !tunnel.IsError || strings.Contains(textOf(tunnel), "SHARED-MARKER") {
			t.Fatalf("tunnel read of extra root was allowed: %s", textOf(tunnel))
		}
		httpServer.Close()

		// The remote server's instructions only mention extra roots when
		// they are usable remotely; local paths are not disclosed otherwise.
		remoteInstructions := s.serverInstructions(workspace.OriginRemote)
		if strings.Contains(remoteInstructions, shared) != remoteEnabled {
			t.Fatalf("remote instructions (remote=%v): %s", remoteEnabled, remoteInstructions)
		}
		if !strings.Contains(s.serverInstructions(workspace.OriginLocal), shared) {
			t.Fatal("local instructions do not list the extra root")
		}
	}
}

func callRead(t *testing.T, endpoint string, header map[string]string, path string) *mcp.CallToolResult {
	t.Helper()
	ctx := context.Background()
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		clone := req.Clone(req.Context())
		for name, value := range header {
			clone.Header.Set(name, value)
		}
		return http.DefaultTransport.RoundTrip(clone)
	})}
	client := mcp.NewClient(&mcp.Implementation{Name: "origin-test", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: endpoint, HTTPClient: httpClient}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	opened, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "open_default_workspace", Arguments: map[string]any{}})
	if err != nil || opened.IsError {
		t.Fatalf("open_default_workspace: %v %v", err, opened)
	}
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "read",
		Arguments: map[string]any{"workspaceId": "default", "path": path},
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
