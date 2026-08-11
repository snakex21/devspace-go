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
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/locales"
	"github.com/snakex21/devspace-go/internal/logger"
	"github.com/snakex21/devspace-go/internal/store"
	"github.com/snakex21/devspace-go/internal/tools"
	"github.com/snakex21/devspace-go/internal/workspace"
)

// boolPtr returns a pointer to the given bool value (for ToolAnnotations pointer fields).
func boolPtr(b bool) *bool { return &b }

const cloudflaredURLTimeout = 30 * time.Second

// Server represents the running MCP WebCoder server.
type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	tunnelStop context.CancelFunc
	registry   *workspace.Registry
	store      *store.Store
}

// New creates a new MCP WebCoder server.
func New(cfg *config.Config) (*Server, error) {
	logger.Init(string(cfg.Logging.Level), string(cfg.Logging.Format))
	tools.SetShell(cfg.Shell)

	s, err := store.New(cfg.StateDir)
	if err != nil {
		return nil, fmt.Errorf("init store: %w", err)
	}

	registry := workspace.NewRegistry(cfg, s)

	return &Server{
		cfg:      cfg,
		registry: registry,
		store:    s,
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

	// MCP endpoint using stateless Streamable HTTP. Workspace state is tracked
	// separately by workspaceId, so transport sessions only add another failure
	// mode when a proxy or web client drops an Mcp-Session-Id between requests.
	handler := s.streamableMCPHandler()

	mux.Handle("/mcp", handler)

	// Legacy SSE endpoint for MCP clients that still expect /sse.
	sseHandler := mcp.NewSSEHandler(
		func(r *http.Request) *mcp.Server {
			return s.createMcpServer()
		},
		&mcp.SSEOptions{DisableLocalhostProtection: true},
	)
	mux.Handle("/sse", sseHandler)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler: s.loggingMiddleware(mux),
	}

	// Auto-start Cloudflare Tunnel if available
	s.startTunnel() // non-fatal

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msg(locales.T("server.shutdown"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown error")
		}
		if s.tunnelStop != nil {
			s.tunnelStop()
		}
		if s.store != nil {
			s.store.Close()
		}
		close(idleConnsClosed)
	}()

	log.Info().
		Str("host", s.cfg.Host).
		Int("port", s.cfg.Port).
		Msg(locales.T("server.listening"))

	log.Info().
		Strs("allowed_roots", s.cfg.AllowedRoots).
		Msg(locales.T("server.roots"))

	log.Info().
		Bool("skills", s.cfg.SkillsEnabled).
		Str("tool_mode", string(s.cfg.ToolMode)).
		Str("tool_naming", string(s.cfg.ToolNaming)).
		Msg(locales.T("server.config"))

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	<-idleConnsClosed
	return nil
}

func (s *Server) streamableMCPHandler() http.Handler {
	return mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return s.createMcpServer()
		},
		&mcp.StreamableHTTPOptions{
			Stateless:                  true,
			JSONResponse:               true,
			DisableLocalhostProtection: true,
		},
	)
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
	fmt.Printf("🔗  %s\n", locales.T("tunnel.starting_pinggy"))
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
		fmt.Printf("⚠️  %s (pinggy): %v\n", locales.T("error.cmd_failed"), err)
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
		s.tunnelStop = cancel
		printTunnelURL(url)
		return url
	case <-time.After(15 * time.Second):
		fmt.Printf("⚠️  %s\n", locales.T("tunnel.pinggy_timeout"))
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
	fmt.Printf("🔗  %s\n", locales.T("tunnel.starting_cloudflared"))
	fmt.Printf("    %s\n", tunnelExe)
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, tunnelExe, "tunnel", "--url", fmt.Sprintf("http://%s:%d", s.cfg.Host, s.cfg.Port))

	stdout, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("⚠️  %s (cloudflared): %v\n", locales.T("error.cmd_failed"), err)
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
		s.tunnelStop = cancel
		printTunnelURL(url)
		return url
	case <-time.After(cloudflaredURLTimeout):
		fmt.Printf("⚠️  %s\n", locales.T("tunnel.cloudflared_timeout"))
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
	lines := []string{
		"🌐 " + locales.T("tunnel.active"),
		mcpURL,
		sseURL,
		"",
		locales.T("tunnel.paste_chatgpt"),
		mcpURL,
		locales.T("tunnel.try_sse"),
	}
	width := 54
	for _, line := range lines {
		if lineWidth := utf8.RuneCountInString(line); lineWidth > width {
			width = lineWidth
		}
	}

	fmt.Println()
	fmt.Printf("╔%s╗\n", strings.Repeat("═", width+4))
	for _, line := range lines {
		padding := strings.Repeat(" ", width-utf8.RuneCountInString(line))
		fmt.Printf("║  %s%s  ║\n", line, padding)
	}
	fmt.Printf("╚%s╝\n", strings.Repeat("═", width+4))
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

	// mkdir
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        names.Mkdir,
			Description: "Create a directory inside an open workspace, including missing parent directories. Use this instead of shell mkdir/New-Item. Call open_workspace first and pass workspaceId.",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true,
				OpenWorldHint:   boolPtr(false),
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input tools.MkdirInput) (*mcp.CallToolResult, tools.MkdirOutput, error) {
			ws, err := s.registry.GetWorkspace(input.WorkspaceID)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.MkdirOutput{}, nil
			}

			_, err = s.registry.ResolvePath(ws, input.Path)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.MkdirOutput{}, nil
			}

			return tools.MakeDirectory(ctx, req, input, ws.Root)
		},
	)

	// move
	mcp.AddTool(server,
		&mcp.Tool{
			Name:        names.Move,
			Description: "Move or rename a file/directory inside an open workspace. Creates missing parent directories for the destination. Use this instead of shell Move-Item/mv. Call open_workspace first and pass workspaceId.",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				DestructiveHint: boolPtr(true),
				IdempotentHint:  false,
				OpenWorldHint:   boolPtr(false),
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input tools.MoveInput) (*mcp.CallToolResult, tools.MoveOutput, error) {
			ws, err := s.registry.GetWorkspace(input.WorkspaceID)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.MoveOutput{}, nil
			}

			_, err = s.registry.ResolvePath(ws, input.SourcePath)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.MoveOutput{}, nil
			}
			_, err = s.registry.ResolvePath(ws, input.TargetPath)
			if err != nil {
				result := &mcp.CallToolResult{}
				result.SetError(err)
				return result, tools.MoveOutput{}, nil
			}

			return tools.MovePath(ctx, req, input, ws.Root)
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
		"Run a shell command inside an open workspace. On Windows, uses PowerShell.exe. On Unix, uses bash. Use only for tests, builds, git inspection, and commands that are better executed by the shell. Do not use %s to create, move, rename, or modify files. Prefer %s for file inspection, %s for creating directories, %s for moves/renames, and %s/%s for file changes. Call open_workspace first and pass workspaceId.",
		names.Bash, names.Read, names.Mkdir, names.Move, names.Edit, names.Write,
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
	Mkdir string
	Move  string
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
			Mkdir: "create_directory",
			Move:  "move_path",
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
		Mkdir: "mkdir",
		Move:  "move",
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
		"Use MCP WebCoder as a local coding workspace. Call open_workspace once per project folder or worktree to obtain a workspaceId; if local absolute paths are blocked by the client, call open_default_workspace instead. Reuse that same workspaceId for all later file, search, edit, write, mkdir, move, and shell tools in that folder. If the workspaceId becomes stale after reconnecting, pass workspaceId 'default' or 'latest' to use the most recent/default workspace. %s%sPrefer %s for targeted modifications, %s only for new files or complete rewrites, %s for directory creation, %s for moves/renames, and %s for tests, builds, git inspection, package scripts, and commands that are better executed by the shell. Do not create, move, rename, or modify files with %s. On Windows, %s uses PowerShell.exe; on Unix, bash.",
		agentsMd,
		inspection,
		names.Edit, names.Write, names.Mkdir, names.Move, names.Bash, names.Bash, names.Bash,
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
