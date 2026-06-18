package server

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
	"github.com/waishnav/mcp-webcoder/internal/auth"
	"github.com/waishnav/mcp-webcoder/internal/config"
	"github.com/waishnav/mcp-webcoder/internal/logger"
	"github.com/waishnav/mcp-webcoder/internal/store"
	"github.com/waishnav/mcp-webcoder/internal/tools"
	"github.com/waishnav/mcp-webcoder/internal/workspace"
)

// boolPtr returns a pointer to the given bool value (for ToolAnnotations pointer fields).
func boolPtr(b bool) *bool { return &b }

// Server represents the running MCP WebCoder server.
type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	provider   *auth.Provider
	registry   *workspace.Registry
	store      *store.Store
	noAuth     bool
}

// New creates a new MCP WebCoder server.
func New(cfg *config.Config, noAuth bool) (*Server, error) {
	logger.Init(string(cfg.Logging.Level), string(cfg.Logging.Format))
	tools.SetShell(cfg.Shell)

	s, err := store.New(cfg.StateDir)
	if err != nil {
		return nil, fmt.Errorf("init store: %w", err)
	}

	registry := workspace.NewRegistry(cfg, s)
	provider := auth.NewProvider(cfg)

	return &Server{
		cfg:      cfg,
		provider: provider,
		registry: registry,
		store:    s,
		noAuth:   noAuth,
	}, nil
}

// Start begins listening for connections.
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":true,"name":"mcp-webcoder"}`)
	})

	// OAuth endpoints
	mux.HandleFunc("/.well-known/oauth-authorization-server", s.provider.HandleOAuthMetadata)
	mux.HandleFunc("/.well-known/oauth-protected-resource/mcp", s.provider.HandleProtectedResourceMetadata)
	mux.HandleFunc("/authorize", s.provider.HandleAuthorize)
	mux.HandleFunc("/token", s.provider.HandleToken)
	mux.HandleFunc("/revoke", s.provider.HandleRevoke)
	mux.HandleFunc("/register", s.provider.HandleRegister)

	// MCP endpoint using StreamableHTTPHandler
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return s.createMcpServer()
		},
		&mcp.StreamableHTTPOptions{DisableLocalhostProtection: true},
	)

	mux.Handle("/mcp", handler)

	// Legacy SSE endpoint for MCP clients that still expect /sse.
	sseHandler := mcp.NewSSEHandler(
		func(r *http.Request) *mcp.Server {
			return s.createMcpServer()
		},
		&mcp.SSEOptions{DisableLocalhostProtection: true},
	)
	mux.Handle("/sse", sseHandler)

	// Wrap with auth middleware (skip if --no-auth)
	var finalHandler http.Handler
	if s.noAuth {
		fmt.Println("⚠️  UWAGA: Autoryzacja wyłączona (--no-auth). Każdy z dostępem do URL może używać serwera.")
		finalHandler = s.loggingMiddleware(mux)
	} else {
		authHandler := s.provider.AuthMiddleware(mux)
		finalHandler = s.loggingMiddleware(authHandler)
	}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler: finalHandler,
	}

	// Auto-start Cloudflare Tunnel if available
	s.startTunnel() // non-fatal

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msg("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown error")
		}
		if s.store != nil {
			s.store.Close()
		}
		close(idleConnsClosed)
	}()

	log.Info().
		Str("host", s.cfg.Host).
		Int("port", s.cfg.Port).
		Msg("webcoder listening")

	log.Info().
		Strs("allowed_roots", s.cfg.AllowedRoots).
		Msg("allowed roots")

	log.Info().
		Bool("skills", s.cfg.SkillsEnabled).
		Str("tool_mode", string(s.cfg.ToolMode)).
		Str("tool_naming", string(s.cfg.ToolNaming)).
		Msg("configuration")

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	<-idleConnsClosed
	return nil
}

// startTunnel attempts to start a tunnel to expose the server publicly.
// Tries cloudflared first, falls back to pinggy.
// Returns the public URL if successful. Non-fatal.
func (s *Server) startTunnel() string {
	// Try cloudflared first
	if url := s.startCloudflared(); url != "" {
		return url
	}

	// Fallback to pinggy
	if url := s.startPinggy(); url != "" {
		return url
	}

	return ""
}

// startPinggy creates a tunnel via pinggy.io using SSH.
// Uses the same SSH key each time → same URL across restarts.
func (s *Server) startPinggy() string {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return "" // ssh not available
	}

	fmt.Println()
	fmt.Println("🔗  Uruchamiam tunel pinggy (stały URL)...")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	target := fmt.Sprintf("R0:localhost:%d", s.cfg.Port)
	cmd := exec.CommandContext(ctx, sshPath,
		"-p", "443",
		"-o", "StrictHostKeyChecking=accept-new",
		"-o", "ServerAliveInterval=30",
		"-R", target,
		"a.pinggy.io",
	)

	stdout, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("⚠️  Nie można uruchomić pinggy: %v\n", err)
		cancel()
		return ""
	}

	// Pinggy prints URL to stdout
	urlRegex := regexp.MustCompile(`https://[a-zA-Z0-9]+\.(a\.)?pinggy\.(link|io|xyz)`)
	done := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if match := urlRegex.FindString(line); match != "" {
				done <- match
				return
			}
		}
	}()

	// Also check stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if match := urlRegex.FindString(line); match != "" {
				select {
				case done <- match:
				default:
				}
				return
			}
		}
	}()

	select {
	case url := <-done:
		printTunnelURL(url)
		return url
	case <-time.After(15 * time.Second):
		fmt.Println("⚠️  pinggy nie zwrócił URL — próbuję cloudflared...")
		cancel()
		return ""
	}
}

// startCloudflared creates a tunnel via cloudflared.
func (s *Server) startCloudflared() string {
	tunnelExe := findCloudflaredExecutable()
	if tunnelExe == "" {
		return ""
	}

	fmt.Println()
	fmt.Println("🔗  Uruchamiam tunel cloudflared...")
	fmt.Printf("    %s\n", tunnelExe)
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, tunnelExe, "tunnel", "--url", fmt.Sprintf("http://%s:%d", s.cfg.Host, s.cfg.Port))

	stdout, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("⚠️  Nie można uruchomić cloudflared: %v\n", err)
		cancel()
		return ""
	}

	urlRegex := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)
	done := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if match := urlRegex.FindString(scanner.Text()); match != "" {
				done <- match
				return
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if match := urlRegex.FindString(line); match != "" {
				select {
				case done <- match:
				default:
				}
				return
			}
		}
	}()

	select {
	case url := <-done:
		printTunnelURL(url)
		return url
	case <-time.After(10 * time.Second):
		fmt.Println("⚠️  cloudflared nie zwrócił URL — serwer działa tylko lokalnie")
		cancel()
		return ""
	}
}

func findCloudflaredExecutable() string {
	names := []string{"cloudflared.exe", "cloudflared"}
	var dirs []string

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		for dir := exeDir; dir != ""; dir = filepath.Dir(dir) {
			dirs = append(dirs, filepath.Join(dir, "tools"), dir)
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(wd, "tools"), wd)
	}

	seen := map[string]bool{}
	for _, dir := range dirs {
		cleanDir := filepath.Clean(dir)
		if seen[cleanDir] {
			continue
		}
		seen[cleanDir] = true
		for _, name := range names {
			candidate := filepath.Join(cleanDir, name)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
	}

	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}

func printTunnelURL(url string) {
	mcpURL := url + "/mcp"
	sseURL := url + "/sse"
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║  🌐 TUNEL AKTYWNY                                       ║")
	fmt.Printf("║  %-54s ║\n", mcpURL)
	fmt.Printf("║  %-54s ║\n", sseURL)
	fmt.Println("║                                                          ║")
	fmt.Println("║  Wklej w ChatGPT jako adres serwera MCP:                 ║")
	fmt.Printf("║  %-54s ║\n", mcpURL)
	fmt.Println("║  Jeśli ChatGPT marudzi, użyj wersji /sse.                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// createMcpServer creates a new MCP server with all tools registered.
func (s *Server) createMcpServer() *mcp.Server {
	mcpServer := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-webcoder", Version: "0.1.0"},
		&mcp.ServerOptions{
			Instructions: s.serverInstructions(),
		},
	)

	s.registerTools(mcpServer)
	return mcpServer
}

// registerTools registers all MCP WebCoder tools on the MCP server.
func (s *Server) registerTools(server *mcp.Server) {
	names := s.toolNames()

	// open_workspace
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "open_workspace",
			Description: "Open a local project directory as a coding workspace. If path is empty or 'default', opens the first configured allowed root. Call this once per project folder or worktree before reading, editing, searching, writing, or running commands. Reuse the returned workspaceId for later calls in the same folder. If a remote client blocks local absolute paths, call open_default_workspace instead.",
			Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input OpenWorkspaceInput) (*mcp.CallToolResult, OpenWorkspaceOutput, error) {
			mode := workspace.ModeCheckout
			if input.Mode == "worktree" {
				mode = workspace.ModeWorktree
			}

			wsCtx, err := s.registry.OpenWorkspace(input.Path, mode, input.BaseRef)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(fmt.Errorf("failed to open workspace: %v", err))
				return result, OpenWorkspaceOutput{}, nil
			}

			var agentsFiles []AgentsFileOutput
			for _, f := range wsCtx.AgentsFiles {
				agentsFiles = append(agentsFiles, AgentsFileOutput{
					Path:    workspace.FormatPath(f.Path, wsCtx.Workspace.Root),
					Content: f.Content,
				})
			}
			var availableAgentsFiles []AvailableAgentsFileOutput
			for _, f := range wsCtx.AvailableAgentsFiles {
				availableAgentsFiles = append(availableAgentsFiles, AvailableAgentsFileOutput{
					Path: workspace.FormatPath(f.Path, wsCtx.Workspace.Root),
				})
			}

			instruction := "Use this workspaceId in all subsequent tool calls for this project. Do not call open_workspace again for this same folder unless this workspaceId stops working, the user asks to reopen, or you switch to a different folder/worktree."

			resultText := fmt.Sprintf("Opened workspace %s\nRoot: %s\nMode: %s\n%s",
				wsCtx.Workspace.ID, wsCtx.Workspace.Root, wsCtx.Workspace.Mode, instruction)

			return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{Text: resultText}},
				}, OpenWorkspaceOutput{
					WorkspaceID:          wsCtx.Workspace.ID,
					Root:                 wsCtx.Workspace.Root,
					Mode:                 string(wsCtx.Workspace.Mode),
					AgentsFiles:          agentsFiles,
					AvailableAgentsFiles: availableAgentsFiles,
					Instruction:          instruction,
				}, nil
		},
	)

	// open_default_workspace avoids passing local absolute paths through clients
	// that may block filesystem-looking arguments before they reach MCP WebCoder.
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "open_default_workspace",
			Description: "Open the default configured workspace without sending a local path. Use this when open_workspace with an absolute Windows/macOS/Linux path is blocked by the MCP client. Returns a workspaceId for the first allowed root.",
			Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input OpenDefaultWorkspaceInput) (*mcp.CallToolResult, OpenWorkspaceOutput, error) {
			wsCtx, err := s.registry.OpenDefaultWorkspace()
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(fmt.Errorf("failed to open default workspace: %v", err))
				return result, OpenWorkspaceOutput{}, nil
			}

			var agentsFiles []AgentsFileOutput
			for _, f := range wsCtx.AgentsFiles {
				agentsFiles = append(agentsFiles, AgentsFileOutput{
					Path:    workspace.FormatPath(f.Path, wsCtx.Workspace.Root),
					Content: f.Content,
				})
			}
			var availableAgentsFiles []AvailableAgentsFileOutput
			for _, f := range wsCtx.AvailableAgentsFiles {
				availableAgentsFiles = append(availableAgentsFiles, AvailableAgentsFileOutput{
					Path: workspace.FormatPath(f.Path, wsCtx.Workspace.Root),
				})
			}

			instruction := "Use this workspaceId in all subsequent tool calls for this project. You may also pass workspaceId 'default' or 'latest' if the exact ID is stale after reconnecting."
			resultText := fmt.Sprintf("Opened default workspace %s\nRoot: %s\nMode: %s\n%s",
				wsCtx.Workspace.ID, wsCtx.Workspace.Root, wsCtx.Workspace.Mode, instruction)

			return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{Text: resultText}},
				}, OpenWorkspaceOutput{
					WorkspaceID:          wsCtx.Workspace.ID,
					Root:                 wsCtx.Workspace.Root,
					Mode:                 string(wsCtx.Workspace.Mode),
					AgentsFiles:          agentsFiles,
					AvailableAgentsFiles: availableAgentsFiles,
					Instruction:          instruction,
				}, nil
		},
	)

	// read
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        names.Read,
			Description: "Read a file inside an open workspace. Use this for file inspection instead of shell commands like cat. Call open_workspace first and pass workspaceId.",
			Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input tools.ReadInput) (*mcp.CallToolResult, tools.ReadOutput, error) {
			ws, err := s.registry.GetWorkspace(input.WorkspaceID)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.ReadOutput{}, nil
			}

			_, err = s.registry.ResolvePath(ws, input.Path)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.ReadOutput{}, nil
			}

			return tools.ReadFile(ctx, req, input, ws.Root)
		},
	)

	// write
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        names.Write,
			Description: fmt.Sprintf("Create or completely overwrite a file inside an open workspace. Prefer %s for targeted changes to existing files. Call open_workspace first and pass workspaceId.", names.Edit),
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				DestructiveHint: boolPtr(true),
				IdempotentHint:  false,
				OpenWorldHint:   boolPtr(false),
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input tools.WriteInput) (*mcp.CallToolResult, tools.WriteOutput, error) {
			ws, err := s.registry.GetWorkspace(input.WorkspaceID)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.WriteOutput{}, nil
			}

			_, err = s.registry.ResolvePath(ws, input.Path)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.WriteOutput{}, nil
			}

			return tools.WriteFile(ctx, req, input, ws.Root)
		},
	)

	// edit
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        names.Edit,
			Description: fmt.Sprintf("Edit one file inside an open workspace by replacing exact text blocks. Prefer this over %s for targeted changes. Call open_workspace first and pass workspaceId.", names.Write),
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: boolPtr(true),
				IdempotentHint:  false,
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input tools.EditInput) (*mcp.CallToolResult, tools.EditOutput, error) {
			ws, err := s.registry.GetWorkspace(input.WorkspaceID)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.EditOutput{}, nil
			}

			_, err = s.registry.ResolvePath(ws, input.Path)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.EditOutput{}, nil
			}

			return tools.EditFile(ctx, req, input, ws.Root)
		},
	)

	// Full mode tools: grep, glob, ls
	if s.cfg.ToolMode == config.ToolModeFull {
		// grep
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        names.Grep,
				Description: "Search file contents inside an open workspace. Use this before broad reads when looking for symbols, text, or usage sites. Call open_workspace first and pass workspaceId.",
				Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
			},
			func(ctx context.Context, req *mcp.CallToolRequest, input tools.GrepInput) (*mcp.CallToolResult, tools.GrepOutput, error) {
				ws, err := s.registry.GetWorkspace(input.WorkspaceID)
				if err != nil {
					result := &mcp.CallToolResult{}
					result.SetError(err)
					return result, tools.GrepOutput{}, nil
				}
				return tools.GrepFiles(ctx, req, input, ws.Root)
			},
		)

		// glob
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        names.Glob,
				Description: "Find files by glob pattern inside an open workspace. Use this to discover filenames or narrow file sets before reading. Call open_workspace first and pass workspaceId.",
				Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
			},
			func(ctx context.Context, req *mcp.CallToolRequest, input tools.GlobInput) (*mcp.CallToolResult, tools.GlobOutput, error) {
				ws, err := s.registry.GetWorkspace(input.WorkspaceID)
				if err != nil {
					result := &mcp.CallToolResult{}
					result.SetError(err)
					return result, tools.GlobOutput{}, nil
				}
				return tools.FindFiles(ctx, req, input, ws.Root)
			},
		)

		// ls
		mcp.AddTool(server,
			&mcp.Tool{
				Name:        names.Ls,
				Description: "List a directory inside an open workspace. Use this for directory inspection before reading files. Call open_workspace first and pass workspaceId.",
				Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
			},
			func(ctx context.Context, req *mcp.CallToolRequest, input tools.LsInput) (*mcp.CallToolResult, tools.LsOutput, error) {
				ws, err := s.registry.GetWorkspace(input.WorkspaceID)
				if err != nil {
					result := &mcp.CallToolResult{}
					result.SetError(err)
					return result, tools.LsOutput{}, nil
				}
				return tools.ListDirectory(ctx, req, input, ws.Root)
			},
		)
	}

	// bash (PowerShell on Windows, bash on Unix)
	bashDesc := fmt.Sprintf(
		"Run a shell command inside an open workspace. On Windows, uses PowerShell.exe. On Unix, uses bash. Use only for tests, builds, git inspection, and commands that are better executed by the shell. Do not use %s to create or modify files. Prefer %s for file inspection. Call open_workspace first and pass workspaceId.",
		names.Bash, names.Read,
	)
	if s.cfg.ToolMode == config.ToolModeMinimal {
		bashDesc = fmt.Sprintf(
			"Run a shell command inside an open workspace. On Windows, uses PowerShell.exe. On Unix, uses bash. In minimal tool mode, %s, %s, and %s are disabled; use shell commands for search and directory inspection. Do not use %s to create or modify files. Prefer %s for direct file reads. Call open_workspace first and pass workspaceId.",
			names.Grep, names.Glob, names.Ls, names.Bash, names.Read,
		)
	}

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        names.Bash,
			Description: bashDesc,
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: boolPtr(true),
				OpenWorldHint:   boolPtr(true),
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input tools.BashInput) (*mcp.CallToolResult, tools.BashOutput, error) {
			ws, err := s.registry.GetWorkspace(input.WorkspaceID)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.BashOutput{}, nil
			}
			return tools.RunBash(ctx, req, input, ws.Root)
		},
	)
}

// ToolNames holds the tool naming configuration.
type ToolNames struct {
	Read  string
	Write string
	Edit  string
	Grep  string
	Glob  string
	Ls    string
	Bash  string
}

func (s *Server) toolNames() ToolNames {
	if s.cfg.ToolNaming == config.NamingLegacy {
		return ToolNames{
			Read:  "read_file",
			Write: "write_file",
			Edit:  "edit_file",
			Grep:  "grep_files",
			Glob:  "find_files",
			Ls:    "list_directory",
			Bash:  "run_shell",
		}
	}
	return ToolNames{
		Read:  "read",
		Write: "write",
		Edit:  "edit",
		Grep:  "grep",
		Glob:  "glob",
		Ls:    "ls",
		Bash:  "bash",
	}
}

func (s *Server) serverInstructions() string {
	names := s.toolNames()

	inspection := fmt.Sprintf("Prefer %s, %s, %s, and %s for file inspection. ",
		names.Read, names.Grep, names.Glob, names.Ls)
	if s.cfg.ToolMode == config.ToolModeMinimal {
		inspection = fmt.Sprintf("In minimal tool mode, %s, %s, and %s are disabled; use %s with command-line tools such as grep, rg, find, ls, and tree for search and directory inspection. ",
			names.Grep, names.Glob, names.Ls, names.Bash)
	}

	agentsMd := "Follow instructions returned by open_workspace. Before working under a path listed in availableAgentsFiles, use read to inspect that instruction file and follow it. "

	return fmt.Sprintf(
		"Use MCP WebCoder as a local coding workspace. Call open_workspace once per project folder or worktree to obtain a workspaceId; if local absolute paths are blocked by the client, call open_default_workspace instead. Reuse that same workspaceId for all later file, search, edit, write, and shell tools in that folder. If the workspaceId becomes stale after reconnecting, pass workspaceId 'default' or 'latest' to use the most recent/default workspace. %s%sPrefer %s for targeted modifications, %s only for new files or complete rewrites, and %s for tests, builds, git inspection, package scripts, and commands that are better executed by the shell. Do not create or modify files with %s. On Windows, %s uses PowerShell.exe; on Unix, bash.",
		agentsMd,
		inspection,
		names.Edit, names.Write, names.Bash, names.Bash, names.Bash,
	)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		path := r.URL.Path
		if !s.cfg.Logging.Requests {
			return
		}
		if !s.cfg.Logging.Assets && strings.HasPrefix(path, "/mcp-app-assets") {
			return
		}

		log.Info().
			Str("method", r.Method).
			Str("path", path).
			Str("remote_addr", r.RemoteAddr).
			Dur("duration_ms", time.Since(start)).
			Msg("http_request")
	})
}

// --- types ---

type OpenWorkspaceInput struct {
	Path    string `json:"path"`
	Mode    string `json:"mode,omitempty"`
	BaseRef string `json:"baseRef,omitempty"`
}

type OpenDefaultWorkspaceInput struct{}

type OpenWorkspaceOutput struct {
	WorkspaceID          string                      `json:"workspaceId"`
	Root                 string                      `json:"root"`
	Mode                 string                      `json:"mode"`
	AgentsFiles          []AgentsFileOutput          `json:"agentsFiles"`
	AvailableAgentsFiles []AvailableAgentsFileOutput `json:"availableAgentsFiles"`
	Instruction          string                      `json:"instruction"`
}

type AgentsFileOutput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type AvailableAgentsFileOutput struct {
	Path string `json:"path"`
}
