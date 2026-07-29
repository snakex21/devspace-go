package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/snakex21/devspace-go/internal/config"
)

func TestRefreshTokenRotationCanContinue(t *testing.T) {
	cfg := config.DefaultConfig()
	provider := NewProvider(cfg)

	provider.mu.Lock()
	current := provider.createRefreshToken("chatgpt-client", []string{"webcoder"})
	provider.mu.Unlock()

	for rotation := 1; rotation <= 3; rotation++ {
		oldToken := current.Token
		access, replacement, err := provider.RefreshAccessToken(oldToken)
		if err != nil {
			t.Fatalf("rotation %d failed: %v", rotation, err)
		}
		if access == nil || replacement == nil || replacement.Token == "" {
			t.Fatalf("rotation %d did not return both replacement tokens", rotation)
		}
		if _, err := provider.VerifyAccessToken(access.Token); err != nil {
			t.Fatalf("rotation %d returned invalid access token: %v", rotation, err)
		}
		if _, _, err := provider.RefreshAccessToken(oldToken); err == nil {
			t.Fatalf("rotation %d left the old refresh token valid", rotation)
		}
		current = replacement
	}
}

func TestRefreshTokenEndpointReturnsRotatedToken(t *testing.T) {
	cfg := config.DefaultConfig()
	provider := NewProvider(cfg)

	provider.mu.Lock()
	current := provider.createRefreshToken("chatgpt-client", []string{"webcoder"})
	provider.mu.Unlock()

	for rotation := 1; rotation <= 3; rotation++ {
		form := url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {current.Token},
		}
		req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		provider.HandleToken(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("rotation %d returned HTTP %d: %s", rotation, rec.Code, rec.Body.String())
		}

		var response struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("rotation %d returned invalid JSON: %v", rotation, err)
		}
		if response.AccessToken == "" || response.RefreshToken == "" {
			t.Fatalf("rotation %d omitted a replacement token: %s", rotation, rec.Body.String())
		}

		current.Token = response.RefreshToken
	}
}

func TestUnauthorizedResponseAdvertisesOAuthMetadata(t *testing.T) {
	cfg := config.DefaultConfig()
	provider := NewProvider(cfg)
	handler := provider.AuthMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("protected handler ran without a token")
	}))

	req := httptest.NewRequest(http.MethodPost, "http://localhost/mcp", nil)
	req.Host = "quiet-river.trycloudflare.com"
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got HTTP %d, want 401", rec.Code)
	}

	challenge := rec.Header().Get("WWW-Authenticate")
	wantMetadata := `resource_metadata="https://quiet-river.trycloudflare.com/.well-known/oauth-protected-resource/mcp"`
	if !strings.Contains(challenge, wantMetadata) {
		t.Fatalf("challenge %q does not contain %q", challenge, wantMetadata)
	}
	if !strings.Contains(challenge, `scope="webcoder"`) {
		t.Fatalf("challenge %q does not advertise the configured scope", challenge)
	}
}

func TestProtectedResourceMetadataIncludesScopes(t *testing.T) {
	cfg := config.DefaultConfig()
	provider := NewProvider(cfg)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/.well-known/oauth-protected-resource", nil)
	req.Host = "quiet-river.trycloudflare.com"
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	provider.HandleProtectedResourceMetadata(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got HTTP %d, want 200", rec.Code)
	}

	var metadata struct {
		Resource        string   `json:"resource"`
		ScopesSupported []string `json:"scopes_supported"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &metadata); err != nil {
		t.Fatal(err)
	}
	if metadata.Resource != "https://quiet-river.trycloudflare.com/mcp" {
		t.Fatalf("resource = %q", metadata.Resource)
	}
	if len(metadata.ScopesSupported) != 1 || metadata.ScopesSupported[0] != "webcoder" {
		t.Fatalf("scopes_supported = %#v", metadata.ScopesSupported)
	}
}
