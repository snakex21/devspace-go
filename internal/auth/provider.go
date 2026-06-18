package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/waishnav/mcp-webcoder/internal/config"
)

// TokenType represents the type of OAuth token.
type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

// Claims represents JWT claims for access tokens.
type Claims struct {
	jwt.RegisteredClaims
	Scopes []string `json:"scopes"`
}

// TokenInfo holds information about an issued token.
type TokenInfo struct {
	Token     string    `json:"token"`
	Type      TokenType `json:"type"`
	ClientID  string    `json:"clientId"`
	ExpiresAt time.Time `json:"expiresAt"`
	Scopes    []string  `json:"scopes"`
}

// Provider implements a single-user OAuth 2.0 provider for MCP.
type Provider struct {
	cfg        *config.Config
	signingKey []byte

	mu             sync.RWMutex
	authCodes      map[string]*authCodeEntry
	accessTokens   map[string]*TokenInfo
	refreshTokens  map[string]*TokenInfo
	clients        map[string]*clientInfo
}

type authCodeEntry struct {
	ClientID    string
	CodeChallenge string
	ExpiresAt   time.Time
	Scopes      []string
}

type clientInfo struct {
	ClientID     string
	ClientSecret string
	RedirectURIs []string
}

// NewProvider creates a new OAuth provider.
func NewProvider(cfg *config.Config) *Provider {
	// Generate a signing key
	signingKey := make([]byte, 32)
	rand.Read(signingKey)

	p := &Provider{
		cfg:           cfg,
		signingKey:    signingKey,
		authCodes:     make(map[string]*authCodeEntry),
		accessTokens:  make(map[string]*TokenInfo),
		refreshTokens: make(map[string]*TokenInfo),
		clients:       make(map[string]*clientInfo),
	}

	// Start cleanup goroutine
	go p.cleanupLoop()

	return p
}

// cleanupLoop periodically removes expired tokens and codes.
func (p *Provider) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		now := time.Now()

		for code, entry := range p.authCodes {
			if now.After(entry.ExpiresAt) {
				delete(p.authCodes, code)
			}
		}
		for token, info := range p.accessTokens {
			if now.After(info.ExpiresAt) {
				delete(p.accessTokens, token)
			}
		}
		for token, info := range p.refreshTokens {
			if now.After(info.ExpiresAt) {
				delete(p.refreshTokens, token)
			}
		}

		p.mu.Unlock()
	}
}

// RegisterClient registers a dynamic OAuth client.
func (p *Provider) RegisterClient(redirectURIs []string) (*clientInfo, error) {
	clientID := generateToken(32)
	clientSecret := generateToken(48)

	client := &clientInfo{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURIs: redirectURIs,
	}

	p.mu.Lock()
	p.clients[clientID] = client
	p.mu.Unlock()

	return client, nil
}

// CreateAuthorizationCode creates an authorization code for the PKCE flow.
func (p *Provider) CreateAuthorizationCode(clientID, codeChallenge string, scopes []string) (string, error) {
	p.mu.RLock()
	_, ok := p.clients[clientID]
	p.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("unknown client")
	}

	code := generateToken(32)
	entry := &authCodeEntry{
		ClientID:      clientID,
		CodeChallenge: codeChallenge,
		ExpiresAt:     time.Now().Add(10 * time.Minute),
		Scopes:        scopes,
	}

	p.mu.Lock()
	p.authCodes[code] = entry
	p.mu.Unlock()

	return code, nil
}

// ExchangeCode exchanges an authorization code + code_verifier for tokens.
func (p *Provider) ExchangeCode(code, codeVerifier string) (*TokenInfo, *TokenInfo, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	entry, ok := p.authCodes[code]
	if !ok {
		return nil, nil, fmt.Errorf("invalid authorization code")
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(p.authCodes, code)
		return nil, nil, fmt.Errorf("authorization code expired")
	}

	// Verify PKCE
	challenge := computeCodeChallenge(codeVerifier)
	if subtle.ConstantTimeCompare([]byte(challenge), []byte(entry.CodeChallenge)) != 1 {
		return nil, nil, fmt.Errorf("invalid code verifier")
	}

	// Delete one-time code
	delete(p.authCodes, code)

	// Create tokens
	accessToken := p.createAccessToken(entry.ClientID, entry.Scopes)
	refreshToken := p.createRefreshToken(entry.ClientID, entry.Scopes)

	return accessToken, refreshToken, nil
}

// RefreshAccessToken creates a new access token from a refresh token.
func (p *Provider) RefreshAccessToken(refreshTokenStr string) (*TokenInfo, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	info, ok := p.refreshTokens[refreshTokenStr]
	if !ok {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(info.ExpiresAt) {
		delete(p.refreshTokens, refreshTokenStr)
		return nil, fmt.Errorf("refresh token expired")
	}

	// Rotate refresh token
	delete(p.refreshTokens, refreshTokenStr)

	newAccessToken := p.createAccessToken(info.ClientID, info.Scopes)
	// Create new refresh token for next rotation
	p.createRefreshToken(info.ClientID, info.Scopes)

	return newAccessToken, nil
}

// VerifyAccessToken validates an access token.
func (p *Provider) VerifyAccessToken(tokenStr string) (*Claims, error) {
	// Try JWT verification
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.signingKey, nil
	})

	if err != nil {
		// Fall back to opaque token lookup
		p.mu.RLock()
		info, ok := p.accessTokens[tokenStr]
		p.mu.RUnlock()

		if !ok {
			return nil, fmt.Errorf("invalid access token")
		}
		if time.Now().After(info.ExpiresAt) {
			return nil, fmt.Errorf("access token expired")
		}

		return &Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:  info.ClientID,
				IssuedAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			},
			Scopes: info.Scopes,
		}, nil
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetOwnerPassword returns the configured owner token hash for verification.
func (p *Provider) VerifyOwnerPassword(password string) bool {
	expectedHash := sha256Hash(p.cfg.OAuth.OwnerToken)
	providedHash := sha256Hash(password)
	return subtle.ConstantTimeCompare([]byte(expectedHash), []byte(providedHash)) == 1
}

// HandleAuthorize handles the OAuth authorization endpoint.
func (p *Provider) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		q := r.URL.Query()
		clientID := q.Get("client_id")
		redirectURI := q.Get("redirect_uri")
		codeChallenge := q.Get("code_challenge")
		state := q.Get("state")
		scope := q.Get("scope")

		if clientID == "" || redirectURI == "" || codeChallenge == "" {
			http.Error(w, "Missing required OAuth parameters", http.StatusBadRequest)
			return
		}

		// If no owner password is configured, auto-approve OAuth.
		if p.cfg.OAuth.OwnerToken == "" {
			code, err := p.CreateAuthorizationCode(clientID, codeChallenge, scopesFromString(scope, p.cfg.OAuth.Scopes))
			if err != nil {
				http.Error(w, "Failed to create authorization code", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, redirectWithCode(redirectURI, code, state), http.StatusFound)
			return
		}

		// Show authorization form, preserving all OAuth parameters for POST.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, authorizeHTML,
			html.EscapeString(clientID),
			html.EscapeString(redirectURI),
			html.EscapeString(codeChallenge),
			html.EscapeString(state),
			html.EscapeString(scope),
		)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	codeChallenge := r.FormValue("code_challenge")
	state := r.FormValue("state")
	scope := r.FormValue("scope")

	if clientID == "" || redirectURI == "" || codeChallenge == "" {
		http.Error(w, "Missing OAuth parameters", http.StatusBadRequest)
		return
	}

	if !p.VerifyOwnerPassword(password) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, authorizeErrorHTML, "Invalid password")
		return
	}

	// Create auth code using the client's PKCE challenge.
	code, err := p.CreateAuthorizationCode(clientID, codeChallenge, scopesFromString(scope, p.cfg.OAuth.Scopes))
	if err != nil {
		http.Error(w, "Failed to create auth code", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectWithCode(redirectURI, code, state), http.StatusFound)
}

// HandleToken handles the OAuth token endpoint.
func (p *Provider) HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}

	grantType := r.FormValue("grant_type")

	switch grantType {
	case "authorization_code":
		code := r.FormValue("code")
		codeVerifier := r.FormValue("code_verifier")

		accessToken, refreshToken, err := p.ExchangeCode(code, codeVerifier)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"access_token":  accessToken.Token,
			"token_type":    "bearer",
			"expires_in":    int(time.Until(accessToken.ExpiresAt).Seconds()),
			"refresh_token": refreshToken.Token,
			"scope":         strings.Join(accessToken.Scopes, " "),
		})

	case "refresh_token":
		refreshTokenStr := r.FormValue("refresh_token")

		accessToken, err := p.RefreshAccessToken(refreshTokenStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"access_token": accessToken.Token,
			"token_type":   "bearer",
			"expires_in":   int(time.Until(accessToken.ExpiresAt).Seconds()),
			"scope":        strings.Join(accessToken.Scopes, " "),
		})

	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "unsupported_grant_type",
		})
	}
}

// HandleRevoke handles token revocation.
func (p *Provider) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}

	token := r.FormValue("token")

	p.mu.Lock()
	delete(p.accessTokens, token)
	delete(p.refreshTokens, token)
	p.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

// HandleProtectedResourceMetadata serves the OAuth protected resource metadata.
func (p *Provider) HandleProtectedResourceMetadata(w http.ResponseWriter, r *http.Request) {
	baseURL := p.baseURLFromRequest(r)

	metadata := map[string]interface{}{
		"resource":                 baseURL + "/mcp",
		"authorization_servers":    []string{baseURL},
		"bearer_methods_supported": []string{"header"},
		"resource_documentation":   baseURL + "/docs",
	}

	writeJSON(w, http.StatusOK, metadata)
}

// HandleOAuthMetadata serves the OAuth authorization server metadata.
func (p *Provider) HandleOAuthMetadata(w http.ResponseWriter, r *http.Request) {
	baseURL := p.baseURLFromRequest(r)

	metadata := map[string]interface{}{
		"issuer":                                           baseURL,
		"authorization_endpoint":                           baseURL + "/authorize",
		"token_endpoint":                                   baseURL + "/token",
		"revocation_endpoint":                              baseURL + "/revoke",
		"registration_endpoint":                            baseURL + "/register",
		"response_types_supported":                         []string{"code"},
		"grant_types_supported":                            []string{"authorization_code", "refresh_token"},
		"code_challenge_methods_supported":                 []string{"S256"},
		"token_endpoint_auth_methods_supported":            []string{"none"},
		"token_endpoint_auth_signing_alg_values_supported": []string{},
		"scopes_supported":                                 p.cfg.OAuth.Scopes,
		"authorization_response_iss_parameter_supported":   true,
	}

	writeJSON(w, http.StatusOK, metadata)
}

// HandleRegister handles dynamic client registration.
func (p *Provider) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RedirectURIs []string `json:"redirect_uris"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}

	client, err := p.RegisterClient(req.RedirectURIs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_error"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"client_id":              client.ClientID,
		"client_secret":          client.ClientSecret,
		"redirect_uris":          client.RedirectURIs,
		"token_endpoint_auth_method": "none",
		"grant_types":            []string{"authorization_code", "refresh_token"},
	})
}

// baseURLFromRequest constructs the public base URL from request headers.
// Uses X-Forwarded-Proto/Host for reverse proxy support (Cloudflare Tunnel, ngrok).
func (p *Provider) baseURLFromRequest(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if fwdProto := r.Header.Get("X-Forwarded-Proto"); fwdProto != "" {
		scheme = fwdProto
	}

	host := r.Host
	if fwdHost := r.Header.Get("X-Forwarded-Host"); fwdHost != "" {
		host = fwdHost
	}

	return scheme + "://" + host
}

func (p *Provider) createAccessToken(clientID string, scopes []string) *TokenInfo {
	now := time.Now()
	expiresAt := now.Add(time.Duration(p.cfg.OAuth.AccessTokenTTLSeconds) * time.Second)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.cfg.PublicBaseURL,
			Subject:   clientID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        generateToken(16),
		},
		Scopes: scopes,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(p.signingKey)

	info := &TokenInfo{
		Token:     tokenStr,
		Type:      AccessToken,
		ClientID:  clientID,
		ExpiresAt: expiresAt,
		Scopes:    scopes,
	}

	p.accessTokens[tokenStr] = info
	return info
}

func (p *Provider) createRefreshToken(clientID string, scopes []string) *TokenInfo {
	now := time.Now()
	expiresAt := now.Add(time.Duration(p.cfg.OAuth.RefreshTokenTTLSeconds) * time.Second)

	tokenStr := generateToken(64)
	info := &TokenInfo{
		Token:     tokenStr,
		Type:      RefreshToken,
		ClientID:  clientID,
		ExpiresAt: expiresAt,
		Scopes:    scopes,
	}

	p.refreshTokens[tokenStr] = info
	return info
}

// ExtractBearerToken extracts the Bearer token from the Authorization header.
func ExtractBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return strings.TrimPrefix(auth, "Bearer "), nil
}

// AuthMiddleware creates HTTP middleware that requires a valid Bearer token.
func (p *Provider) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for OAuth endpoints and health checks
		path := r.URL.Path
		if strings.HasPrefix(path, "/.well-known/") ||
			path == "/authorize" ||
			path == "/token" ||
			path == "/revoke" ||
			path == "/register" ||
			path == "/healthz" ||
			strings.HasPrefix(path, "/mcp-app-assets") {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr, err := ExtractBearerToken(r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"jsonrpc": "2.0",
				"error":   err.Error(),
				"id":      "null",
			})
			return
		}

		claims, err := p.VerifyAccessToken(tokenStr)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"jsonrpc": "2.0",
				"error": map[string]interface{}{
					"code":    -32001,
					"message": "Unauthorized: " + err.Error(),
				},
				"id": nil,
			})
			return
		}

		// Store claims in context
		_ = claims
		next.ServeHTTP(w, r)
	})
}

// Helper functions

func generateToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}

func computeCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func redirectWithCode(redirectURI, code, state string) string {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return redirectURI
	}
	q := u.Query()
	q.Set("code", code)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func scopesFromString(scope string, fallback []string) []string {
	if strings.TrimSpace(scope) == "" {
		return fallback
	}
	return strings.Fields(scope)
}

func sha256Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// HTML templates for the authorization page
const authorizeHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>MCP WebCoder — Authorization</title>
    <style>
        body { font-family: system-ui; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background: #f5f5f5; }
        .card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,.1); max-width: 400px; width: 100%%; }
        h1 { margin: 0 0 .5rem; font-size: 1.5rem; }
        p { color: #666; margin: .5rem 0 1.5rem; }
        input { width: 100%%; padding: .5rem; margin: .5rem 0; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%%; padding: .75rem; background: #0070f3; color: white; border: none; border-radius: 4px; font-size: 1rem; cursor: pointer; }
        button:hover { background: #0051cc; }
    </style>
</head>
<body>
    <div class="card">
        <h1>MCP WebCoder Authorization</h1>
        <p>Enter the owner password to authorize this MCP client.</p>
        <form method="POST">
            <input type="hidden" name="client_id" value="%s">
            <input type="hidden" name="redirect_uri" value="%s">
            <input type="hidden" name="code_challenge" value="%s">
            <input type="hidden" name="state" value="%s">
            <input type="hidden" name="scope" value="%s">
            <input type="password" name="password" placeholder="Owner password" required autofocus>
            <button type="submit">Authorize</button>
        </form>
    </div>
</body>
</html>`

const authorizeErrorHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>MCP WebCoder — Authorization Failed</title>
    <style>
        body { font-family: system-ui; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; background: #f5f5f5; }
        .card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,.1); max-width: 400px; width: 100%%; }
        h1 { margin: 0 0 .5rem; color: #e00; }
        p { color: #666; }
        a { color: #0070f3; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Authorization Failed</h1>
        <p>%s</p>
        <p><a href="javascript:history.back()">Try again</a></p>
    </div>
</body>
</html>`
