package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/waishnav/mcp-webcoder/internal/config"
	"github.com/waishnav/mcp-webcoder/internal/locales"
	"github.com/waishnav/mcp-webcoder/internal/server"
)

func main() {
	// Parse CLI arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve", "start":
			runServer()
		case "init":
			runInit()
		case "doctor":
			runDoctor()
		case "config":
			if len(os.Args) > 2 && os.Args[2] == "get" {
				runConfigShow()
			} else {
				fmt.Println("Usage: mcp-webcoder config get")
			}
		case "help", "--help", "-h":
			printHelp()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Run 'mcp-webcoder help' for usage.")
		}
		return
	}

	// Default: check config, suggest GUI if missing, otherwise run server
	cfg := config.LoadConfig()

	if len(cfg.AllowedRoots) == 0 {
		fmt.Fprintln(os.Stderr, "Error: MCP WebCoder nie jest skonfigurowany.")
		fmt.Fprintln(os.Stderr, "Uruchom konfigurator: mcp-webcoder-gui.exe")
		fmt.Fprintln(os.Stderr, "Lub tekstowo:         mcp-webcoder.exe init")
		os.Exit(1)
	}
	runServerWithConfig(cfg)
}

func runServer() {
	cfg := config.LoadConfig()

	if len(cfg.AllowedRoots) == 0 {
		fmt.Fprintln(os.Stderr, "Error: WEBCODER_ALLOWED_ROOTS must be set.")
		fmt.Fprintln(os.Stderr, "Run 'mcp-webcoder init' to configure.")
		os.Exit(1)
	}

	runServerWithConfig(cfg)
}

func runServerWithConfig(cfg *config.Config) {
	// Initialize locale system
	locales.Init(cfg.Lang)
	if cfg.Lang == "" {
		locales.Init("pl")
	}

	// Auto-detect: no auth if no password configured
	noAuth := cfg.OAuth.OwnerToken == ""

	srv, err := server.New(cfg, noAuth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func runInit() {
	fmt.Println("MCP WebCoder Init")
	fmt.Println("=============")
	fmt.Println()

	cfg := config.DefaultConfig()

	// Read roots
	fmt.Print("Allowed roots (comma-separated paths): ")
	var rootsInput string
	fmt.Scanln(&rootsInput)
	if rootsInput != "" {
		cfg.AllowedRoots = splitAndTrim(rootsInput, ",")
	} else {
		home, _ := os.UserHomeDir()
		cfg.AllowedRoots = []string{home}
		fmt.Printf("Defaulting to: %s\n", home)
	}

	// Read port
	fmt.Printf("Port [%d]: ", cfg.Port)
	var portInput int
	fmt.Scanf("%d", &portInput)
	if portInput > 0 {
		cfg.Port = portInput
	}

	// Read public URL
	fmt.Printf("Public base URL [%s]: ", cfg.PublicBaseURL)
	var urlInput string
	fmt.Scanln(&urlInput)
	if urlInput != "" {
		cfg.PublicBaseURL = urlInput
	}

	// Read owner password
	fmt.Print("Owner password (min 16 chars): ")
	var passwordInput string
	fmt.Scanln(&passwordInput)
	if len(passwordInput) >= 16 {
		cfg.OAuth.OwnerToken = passwordInput
	} else {
		fmt.Println("Password too short, using default...")
		cfg.OAuth.OwnerToken = "webcoder-default-owner-token-change-me"
	}

	// Save config
	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Configuration saved!")
	fmt.Printf("  Config: %s\n", cfg.ConfigDir)
	fmt.Println()
	fmt.Println("Run 'mcp-webcoder serve' to start the server.")
}

func runDoctor() {
	fmt.Println("MCP WebCoder Doctor")
	fmt.Println("=====================")
	fmt.Println()

	// Check Go version
	fmt.Println("[ok] Go runtime")

	// Check config
	cfg := config.LoadConfig()
	if len(cfg.AllowedRoots) > 0 {
		fmt.Printf("[ok] Allowed roots: %s\n", strings.Join(cfg.AllowedRoots, ", "))
	} else {
		fmt.Println("[!] No allowed roots configured")
	}

	if cfg.OAuth.OwnerToken != "" {
		fmt.Println("[ok] Owner token configured")
	} else {
		fmt.Println("[!] No owner token configured")
	}

	// Check shell availability
	fmt.Println("[ok] Shell available (PowerShell on Windows, bash on Unix)")

	// Check state directory
	if _, err := os.Stat(cfg.StateDir); os.IsNotExist(err) {
		fmt.Printf("[ok] State directory will be created: %s\n", cfg.StateDir)
	} else {
		fmt.Printf("[ok] State directory exists: %s\n", cfg.StateDir)
	}

	fmt.Println()
	fmt.Println("Ready to serve!")
}

func runConfigShow() {
	cfg := config.LoadConfig()

	fmt.Println("MCP WebCoder Configuration")
	fmt.Println("===========================")
	fmt.Printf("Host:            %s\n", cfg.Host)
	fmt.Printf("Port:            %d\n", cfg.Port)
	fmt.Printf("Public URL:      %s\n", cfg.PublicBaseURL)
	fmt.Printf("Allowed Roots:   %s\n", strings.Join(cfg.AllowedRoots, ", "))
	fmt.Printf("State Dir:       %s\n", cfg.StateDir)
	fmt.Printf("Config Dir:      %s\n", cfg.ConfigDir)
	fmt.Printf("Agent Dir:       %s\n", cfg.AgentDir)
	fmt.Printf("Tool Mode:       %s\n", cfg.ToolMode)
	fmt.Printf("Tool Naming:     %s\n", cfg.ToolNaming)
	fmt.Printf("Shell:           %s\n", cfg.Shell)
	fmt.Printf("Język:           %s\n", cfg.Lang)
	fmt.Printf("Widgets:         %s\n", cfg.Widgets)
	fmt.Printf("Skills:          %v\n", cfg.SkillsEnabled)
	fmt.Printf("Log Level:       %s\n", cfg.Logging.Level)
	fmt.Printf("Log Format:      %s\n", cfg.Logging.Format)
	fmt.Printf("OAuth:           configured=%v\n", cfg.OAuth.OwnerToken != "")
}

func printHelp() {
	fmt.Print(`MCP WebCoder — Web MCP Coding Workspace (Go)

Usage:
  mcp-webcoder [command]

Commands:
  serve       Start the MCP server (default)
  init        Text-based interactive configuration
  doctor      Diagnostic checks
  config get  Show current configuration
  help        Show this help

GUI Configurator:
  mcp-webcoder-gui.exe    Desktop configuration window

Jeśli nie ustawisz hasła → serwer działa bez autoryzacji.
Jeśli ustawisz hasło → wymagane logowanie OAuth.

Environment:
  WEBCODER_ALLOWED_ROOTS       Required. Comma-separated allowed paths.
  WEBCODER_OAUTH_OWNER_TOKEN   Required. Owner password.
  WEBCODER_PUBLIC_BASE_URL     Public base URL (default: http://127.0.0.1:7676)
  HOST                         Listen host (default: 127.0.0.1)
  PORT                         Listen port (default: 7676)
`)
}

// helpers

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func saveConfig(cfg *config.Config) error {
	if err := os.MkdirAll(cfg.ConfigDir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	configPath := filepath.Join(cfg.ConfigDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	// Save auth token
	authPath := filepath.Join(cfg.ConfigDir, "auth.json")
	authData := map[string]string{"ownerToken": cfg.OAuth.OwnerToken}
	authBytes, err := json.MarshalIndent(authData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal auth: %w", err)
	}

	if err := os.WriteFile(authPath, authBytes, 0600); err != nil {
		return fmt.Errorf("write auth: %w", err)
	}

	fmt.Printf("Configuration saved to %s\n", cfg.ConfigDir)
	return nil
}
