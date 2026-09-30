package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

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

const Version = "2.1.0"

// Server represents the running Dev Space Go server.
type Server struct {
	cfg      *config.Config
	startMu  sync.Mutex
	started  bool
	registry *workspace.Registry
	store    *store.Store
}

// New creates a new Dev Space Go server.
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

// Start serves until Ctrl+C or SIGTERM. Signals are handled even during tunnel startup.
func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return s.StartContext(ctx)
}

// StartContext serves until ctx is canceled. A Server owns its store and may be
// started only once; create a new Server for another run.
func (s *Server) StartContext(ctx context.Context) error {
	return s.start(ctx, s.startTunnel)
}

func (s *Server) start(ctx context.Context, startTunnel func(context.Context, string) *tunnelProcess) error {
	s.startMu.Lock()
	if s.started {
		s.startMu.Unlock()
		return errors.New("server has already been started")
	}
	s.started = true
	s.startMu.Unlock()
	if s.store != nil {
		defer s.store.Close()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprintf(w, `{"ok":true,"name":"devspace-go"}`)
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

	// Bind before launching a tunnel. An occupied port must not expose another
	// process, and the origin must already serve requests during tunnel startup.
	listener, err := net.Listen("tcp", net.JoinHostPort(s.cfg.Host, fmt.Sprint(s.cfg.Port)))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer listener.Close()
	httpServer := &http.Server{Handler: s.loggingMiddleware(mux)}
	serveDone := make(chan error, 1)
	go func() { serveDone <- httpServer.Serve(listener) }()

	tunnelDone := make(chan struct{})
	go func() {
		defer close(tunnelDone)
		tunnel := startTunnel(ctx, tunnelOrigin(listener.Addr().String()))
		if tunnel == nil {
			return
		}
		defer tunnel.stop()
		select {
		case <-ctx.Done():
		case <-tunnel.done:
			if ctx.Err() == nil {
				log.Error().Err(tunnel.err).Msg("tunnel exited; public URL is no longer active; server remains available locally")
			}
		}
	}()

	log.Info().Str("address", listener.Addr().String()).Msg(locales.T("server.listening"))
	log.Info().Strs("allowed_roots", s.cfg.AllowedRoots).Msg(locales.T("server.roots"))
	log.Info().Bool("skills", s.cfg.SkillsEnabled).
		Str("tool_mode", string(s.cfg.ToolMode)).Str("tool_naming", string(s.cfg.ToolNaming)).
		Msg(locales.T("server.config"))

	select {
	case err = <-serveDone:
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
	case <-ctx.Done():
		log.Info().Msg(locales.T("server.shutdown"))
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
			// Shutdown does not forcibly close long-lived SSE or other active requests.
			_ = httpServer.Close()
		}
		shutdownCancel()
		<-serveDone
	}
	cancel()
	_ = httpServer.Close()
	<-tunnelDone
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Wildcard listen addresses are not valid upstream targets on all platforms.
func tunnelOrigin(address string) string {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "http://" + address
	}
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	if host == "::" {
		host = "::1"
	}
	return "http://" + net.JoinHostPort(host, port)
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

// createMcpServer creates a new MCP server with all tools registered.
func (s *Server) createMcpServer() *mcp.Server {
	mcpServer := mcp.NewServer(
		&mcp.Implementation{Name: "devspace-go", Version: Version},
		&mcp.ServerOptions{
			Instructions: s.serverInstructions(),
		},
	)

	s.registerTools(mcpServer)
	return mcpServer
}

// registerTools registers all Dev Space Go tools on the MCP server.
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
	// that may block filesystem-looking arguments before they reach Dev Space Go.
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
		"Use Dev Space Go as a local coding workspace. Call open_workspace once per project folder or worktree to obtain a workspaceId; if local absolute paths are blocked by the client, call open_default_workspace instead. Reuse that same workspaceId for all later file, search, edit, write, mkdir, move, and shell tools in that folder. If the workspaceId becomes stale after reconnecting, pass workspaceId 'default' or 'latest' to use the most recent/default workspace. %s%sPrefer %s for targeted modifications, %s only for new files or complete rewrites, %s for directory creation, %s for moves/renames, and %s for tests, builds, git inspection, package scripts, and commands that are better executed by the shell. Do not create, move, rename, or modify files with %s. On Windows, %s uses PowerShell.exe; on Unix, bash.",
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
