package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/locales"
	"github.com/snakex21/devspace-go/internal/workspace"
)

func captureTunnelOutput(t *testing.T, lang string) string {
	t.Helper()
	locales.Init(lang)
	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = originalStdout })

	printTunnelURL("https://example.trycloudflare.com")
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = originalStdout

	var output bytes.Buffer
	if _, err := output.ReadFrom(reader); err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	return output.String()
}

func TestTunnelOutputUsesSelectedLocale(t *testing.T) {
	tests := []struct {
		lang       string
		expected   []string
		unexpected []string
	}{
		{
			lang:       "en",
			expected:   []string{"TUNNEL ACTIVE", "Paste this in ChatGPT", "try the /sse version"},
			unexpected: []string{"TUNEL AKTYWNY", "Wklej w ChatGPT", "Jeśli ChatGPT"},
		},
		{
			lang:       "de",
			expected:   []string{"TUNNEL AKTIV", "Füge dies in ChatGPT", "versuche die /sse-Version"},
			unexpected: []string{"TUNEL AKTYWNY", "Paste this in ChatGPT"},
		},
	}

	for _, test := range tests {
		t.Run(test.lang, func(t *testing.T) {
			text := captureTunnelOutput(t, test.lang)
			for _, expected := range test.expected {
				if !strings.Contains(text, expected) {
					t.Fatalf("tunnel output does not contain %q:\n%s", expected, text)
				}
			}
			for _, unexpected := range test.unexpected {
				if strings.Contains(text, unexpected) {
					t.Fatalf("tunnel output unexpectedly contains %q:\n%s", unexpected, text)
				}
			}
		})
	}
}

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

func TestStreamableMCPThroughSessionDroppingProxyAndReconnects(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "fixture.txt"), []byte("proxy reconnect test"), 0600); err != nil {
		t.Fatal(err)
	}

	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}

	server := &Server{
		cfg:      cfg,
		registry: workspace.NewRegistry(cfg, nil),
	}
	upstream := httptest.NewServer(server.streamableMCPHandler())
	defer upstream.Close()

	upstreamURL, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}
	proxyHandler := httputil.NewSingleHostReverseProxy(upstreamURL)
	originalDirector := proxyHandler.Director
	proxyHandler.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Del("Mcp-Session-Id")
		req.Close = true
	}

	proxy := httptest.NewServer(proxyHandler)
	defer proxy.Close()

	httpClient := &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}
	defer httpClient.CloseIdleConnections()

	ctx := context.Background()
	for reconnect := 1; reconnect <= 10; reconnect++ {
		client := mcp.NewClient(&mcp.Implementation{
			Name:    "proxy-reconnect-test",
			Version: "1.0.0",
		}, nil)
		session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
			Endpoint:   proxy.URL,
			HTTPClient: httpClient,
		}, nil)
		if err != nil {
			t.Fatalf("reconnect %d failed: %v", reconnect, err)
		}

		for call := 1; call <= 10; call++ {
			result, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name: "read",
				Arguments: map[string]any{
					"workspaceId": "default",
					"path":        "fixture.txt",
				},
			})
			if err != nil {
				session.Close()
				t.Fatalf("reconnect %d, tool call %d failed: %v", reconnect, call, err)
			}
			if result.IsError {
				session.Close()
				t.Fatalf("reconnect %d, tool call %d returned an MCP error: %#v", reconnect, call, result.Content)
			}
		}

		if err := session.Close(); err != nil {
			t.Fatalf("reconnect %d close failed: %v", reconnect, err)
		}
	}
}
