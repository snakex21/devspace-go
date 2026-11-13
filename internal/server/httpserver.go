package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
)

const (
	// readHeaderTimeout bounds slow-header (Slowloris) connections.
	readHeaderTimeout = 10 * time.Second
	// readTimeout bounds reading a whole request, including a large write
	// tool payload over a slow tunnel.
	readTimeout = 2 * time.Minute
	// idleTimeout closes idle keep-alive connections.
	idleTimeout = 2 * time.Minute
	// shutdownTimeout bounds graceful shutdown before connections are closed.
	shutdownTimeout = 10 * time.Second
	// maxRequestBodyBytes caps one MCP request. The largest legitimate
	// payload is a write/edit of a big source file.
	maxRequestBodyBytes = 64 << 20
	maxHeaderBytes      = 1 << 20
)

// newHTTPServer configures timeouts for every phase except writing the
// response. WriteTimeout stays 0 on purpose: /sse streams and shell commands
// (up to 300s) legitimately keep a response open; they end when the client
// disconnects or the server shuts down (baseCtx is canceled).
func newHTTPServer(handler http.Handler, baseCtx context.Context) *http.Server {
	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
		BaseContext:       func(net.Listener) context.Context { return baseCtx },
	}
}

func limitRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxRequestBodyBytes {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		}
		next.ServeHTTP(w, r)
	})
}

// recoverMiddleware turns a handler panic into a logged 500 instead of a
// dropped connection with an unstructured stderr dump.
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if recovered == http.ErrAbortHandler {
					panic(recovered)
				}
				log.Error().Str("path", r.URL.Path).Interface("panic", recovered).
					Str("stack", string(debug.Stack())).Msg("http handler panicked")
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// addTool registers a tool whose handler cannot crash the server. The MCP
// SDK runs handlers on their own goroutines, where an unrecovered panic
// would terminate the whole process, not just the request.
func addTool[In, Out any](server *mcp.Server, tool *mcp.Tool, handler mcp.ToolHandlerFor[In, Out]) {
	mcp.AddTool(server, tool, recoverTool(tool.Name, handler))
}

func recoverTool[In, Out any](name string, handler mcp.ToolHandlerFor[In, Out]) mcp.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input In) (result *mcp.CallToolResult, output Out, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error().Str("tool", name).Interface("panic", recovered).
					Str("stack", string(debug.Stack())).Msg("tool handler panicked")
				var zero Out
				result, output, err = &mcp.CallToolResult{}, zero, nil
				result.SetError(fmt.Errorf("internal error in tool %s: %v", name, recovered))
			}
		}()
		return handler(ctx, req, input)
	}
}

// describeListenError explains the most common startup failure: the port is
// taken, often by a previous Dev Space Go that is still running.
func describeListenError(address string, err error) error {
	if isAddressInUse(err) {
		if probeDevSpace(address) {
			return fmt.Errorf("listen on %s: another Dev Space Go instance is already running on this port; stop it first or change \"port\" in config.json: %w", address, err)
		}
		return fmt.Errorf("listen on %s: the port is already used by another program; close it or change \"port\" in config.json: %w", address, err)
	}
	return fmt.Errorf("listen on %s: %w", address, err)
}

func probeDevSpace(address string) bool {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return false
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	client := &http.Client{Timeout: time.Second}
	defer client.CloseIdleConnections()
	response, err := client.Get("http://" + net.JoinHostPort(host, port) + "/healthz")
	if err != nil {
		return false
	}
	defer response.Body.Close()
	var health struct {
		Name string `json:"name"`
	}
	return json.NewDecoder(http.MaxBytesReader(nil, response.Body, 4096)).Decode(&health) == nil && health.Name == "devspace-go"
}

var errAddressInUse = errors.New("address in use")
