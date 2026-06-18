package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// WidgetMode controls UI widget attachment to tool responses.
type WidgetMode string

const (
	WidgetOff     WidgetMode = "off"
	WidgetChanges WidgetMode = "changes"
	WidgetFull    WidgetMode = "full"
)

// ToolMode controls which tools are available.
type ToolMode string

const (
	ToolModeMinimal ToolMode = "minimal"
	ToolModeFull    ToolMode = "full"
)

// ToolNaming controls tool name format.
type ToolNaming string

const (
	NamingShort  ToolNaming = "short"
	NamingLegacy ToolNaming = "legacy"
)

// LogLevel represents logging verbosity.
type LogLevel string

const (
	LogSilent LogLevel = "silent"
	LogError  LogLevel = "error"
	LogWarn   LogLevel = "warn"
	LogInfo   LogLevel = "info"
	LogDebug  LogLevel = "debug"
)

// LogFormat controls log output format.
type LogFormat string

const (
	LogJSON LogFormat = "json"
	LogText LogFormat = "text"
)

// OAuthConfig holds OAuth-related configuration.
type OAuthConfig struct {
	OwnerToken             string   `json:"ownerToken"`
	AccessTokenTTLSeconds  int      `json:"accessTokenTTLSeconds"`
	RefreshTokenTTLSeconds int      `json:"refreshTokenTTLSeconds"`
	Scopes                 []string `json:"scopes"`
	AllowedRedirectHosts   []string `json:"allowedRedirectHosts"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level         LogLevel  `json:"level"`
	Format        LogFormat `json:"format"`
	Requests      bool      `json:"requests"`
	Assets        bool      `json:"assets"`
	ToolCalls     bool      `json:"toolCalls"`
	ShellCommands bool      `json:"shellCommands"`
	TrustProxy    bool      `json:"trustProxy"`
}

// Config holds all MCP WebCoder server configuration.
type Config struct {
	Host          string        `json:"host"`
	Port          int           `json:"port"`
	AllowedRoots  []string      `json:"allowedRoots"`
	PublicBaseURL string        `json:"publicBaseUrl"`
	StateDir      string        `json:"stateDir"`
	WorktreeRoot  string        `json:"worktreeRoot"`
	AgentDir      string        `json:"agentDir"`
	ConfigDir     string        `json:"configDir"`
	ToolMode      ToolMode      `json:"toolMode"`
	ToolNaming    ToolNaming    `json:"toolNaming"`
	Shell         string        `json:"shell"`
	Lang          string        `json:"lang"`
	Widgets       WidgetMode    `json:"widgets"`
	SkillsEnabled bool          `json:"skillsEnabled"`
	SkillPaths    []string      `json:"skillPaths"`
	AllowedHosts  []string      `json:"allowedHosts"`
	OAuth         OAuthConfig   `json:"oauth"`
	Logging       LoggingConfig `json:"logging"`
}

// DefaultConfig returns a Config with sensible defaults.
// All paths default to the directory containing the executable (portable mode).
func DefaultConfig() *Config {
	exeDir := exeDir()

	return &Config{
		Host:          "127.0.0.1",
		Port:          7676,
		AllowedRoots:  []string{},
		PublicBaseURL: "http://127.0.0.1:7676",
		StateDir:      filepath.Join(exeDir, ".webcoder-state"),
		WorktreeRoot:  filepath.Join(exeDir, ".webcoder", "worktrees"),
		AgentDir:      filepath.Join(exeDir, ".codex"),
		ConfigDir:     filepath.Join(exeDir, ".webcoder"),
		ToolMode:      ToolModeFull,
		ToolNaming:    NamingShort,
		Shell:         "auto",
		Lang:          "auto",
		Widgets:       WidgetFull,
		SkillsEnabled: true,
		AllowedHosts:  []string{"*"},
		OAuth: OAuthConfig{
			AccessTokenTTLSeconds:  3600,
			RefreshTokenTTLSeconds: 2592000,
			Scopes:                 []string{"webcoder"},
			AllowedRedirectHosts:   []string{"chatgpt.com", "localhost", "127.0.0.1"},
		},
		Logging: LoggingConfig{
			Level:         LogInfo,
			Format:        LogJSON,
			Requests:      true,
			Assets:        false,
			ToolCalls:     true,
			ShellCommands: false,
			TrustProxy:    false,
		},
	}
}

// LoadConfig loads configuration from environment variables and config files.
func LoadConfig() *Config {
	cfg := DefaultConfig()

	// Environment variable overrides
	if v := os.Getenv("HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Port)
	}
	if v := os.Getenv("WEBCODER_ALLOWED_ROOTS"); v != "" {
		cfg.AllowedRoots = splitAndTrim(v, ",")
	}
	if v := os.Getenv("WEBCODER_PUBLIC_BASE_URL"); v != "" {
		cfg.PublicBaseURL = v
	}
	if v := os.Getenv("WEBCODER_STATE_DIR"); v != "" {
		cfg.StateDir = v
	}
	if v := os.Getenv("WEBCODER_WORKTREE_ROOT"); v != "" {
		cfg.WorktreeRoot = v
	}
	if v := os.Getenv("WEBCODER_AGENT_DIR"); v != "" {
		cfg.AgentDir = v
	}
	if v := os.Getenv("WEBCODER_CONFIG_DIR"); v != "" {
		cfg.ConfigDir = v
	}
	if v := os.Getenv("WEBCODER_TOOL_MODE"); v != "" {
		cfg.ToolMode = ToolMode(v)
	}
	if v := os.Getenv("WEBCODER_TOOL_NAMING"); v != "" {
		cfg.ToolNaming = ToolNaming(v)
	}
	if v := os.Getenv("WEBCODER_SHELL"); v != "" {
		cfg.Shell = v
	}
	if v := os.Getenv("WEBCODER_LANG"); v != "" {
		cfg.Lang = v
	}
	// Resolve "auto" to system language
	if cfg.Lang == "auto" || cfg.Lang == "" {
		if detected := detectSystemLang(); detected != "" {
			cfg.Lang = detected
		} else {
			cfg.Lang = "en"
		}
	}
	if v := os.Getenv("WEBCODER_WIDGETS"); v != "" {
		cfg.Widgets = WidgetMode(v)
	}
	if v := os.Getenv("WEBCODER_SKILLS"); v == "0" {
		cfg.SkillsEnabled = false
	}
	if v := os.Getenv("WEBCODER_SKILL_PATHS"); v != "" {
		cfg.SkillPaths = splitAndTrim(v, ",")
	}
	if v := os.Getenv("WEBCODER_ALLOWED_HOSTS"); v != "" {
		cfg.AllowedHosts = splitAndTrim(v, ",")
	}

	// OAuth config
	if v := os.Getenv("WEBCODER_OAUTH_OWNER_TOKEN"); v != "" {
		cfg.OAuth.OwnerToken = v
	}
	if v := os.Getenv("WEBCODER_OAUTH_ACCESS_TOKEN_TTL_SECONDS"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.OAuth.AccessTokenTTLSeconds)
	}
	if v := os.Getenv("WEBCODER_OAUTH_REFRESH_TOKEN_TTL_SECONDS"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.OAuth.RefreshTokenTTLSeconds)
	}
	if v := os.Getenv("WEBCODER_OAUTH_SCOPES"); v != "" {
		cfg.OAuth.Scopes = splitAndTrim(v, ",")
	}
	if v := os.Getenv("WEBCODER_OAUTH_ALLOWED_REDIRECT_HOSTS"); v != "" {
		cfg.OAuth.AllowedRedirectHosts = splitAndTrim(v, ",")
	}

	// Logging config
	if v := os.Getenv("WEBCODER_LOG_LEVEL"); v != "" {
		cfg.Logging.Level = LogLevel(v)
	}
	if v := os.Getenv("WEBCODER_LOG_FORMAT"); v != "" {
		cfg.Logging.Format = LogFormat(v)
	}
	if v := os.Getenv("WEBCODER_LOG_REQUESTS"); v == "0" {
		cfg.Logging.Requests = false
	}
	if v := os.Getenv("WEBCODER_LOG_ASSETS"); v == "1" {
		cfg.Logging.Assets = true
	}
	if v := os.Getenv("WEBCODER_LOG_TOOL_CALLS"); v == "0" {
		cfg.Logging.ToolCalls = false
	}
	if v := os.Getenv("WEBCODER_LOG_SHELL_COMMANDS"); v == "1" {
		cfg.Logging.ShellCommands = true
	}
	if v := os.Getenv("WEBCODER_TRUST_PROXY"); v == "1" {
		cfg.Logging.TrustProxy = true
	}

	// Load from config file if exists (new path first, old path as migration fallback)
	configFile := filepath.Join(cfg.ConfigDir, "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Migration: try old .devspace/config.json
		oldConfigFile := filepath.Join(filepath.Dir(cfg.ConfigDir), ".devspace", "config.json")
		if _, err := os.Stat(oldConfigFile); err == nil {
			configFile = oldConfigFile
		}
	}
	if data, err := os.ReadFile(configFile); err == nil {
		data = stripBOM(data)
		var fileConfig Config
		var raw map[string]json.RawMessage
		_ = json.Unmarshal(data, &raw)
		if err := json.Unmarshal(data, &fileConfig); err == nil {
			// Merge file config (file takes precedence for non-empty values)
			if fileConfig.Host != "" {
				cfg.Host = fileConfig.Host
			}
			if fileConfig.Port != 0 {
				cfg.Port = fileConfig.Port
			}
			if len(fileConfig.AllowedRoots) > 0 {
				cfg.AllowedRoots = fileConfig.AllowedRoots
			}
			if fileConfig.PublicBaseURL != "" {
				cfg.PublicBaseURL = fileConfig.PublicBaseURL
			}
			if fileConfig.StateDir != "" {
				cfg.StateDir = fileConfig.StateDir
			}
			if fileConfig.WorktreeRoot != "" {
				cfg.WorktreeRoot = fileConfig.WorktreeRoot
			}
			if fileConfig.AgentDir != "" {
				cfg.AgentDir = fileConfig.AgentDir
			}
			if fileConfig.ConfigDir != "" {
				cfg.ConfigDir = fileConfig.ConfigDir
			}
			if fileConfig.ToolMode != "" {
				cfg.ToolMode = fileConfig.ToolMode
			}
			if fileConfig.ToolNaming != "" {
				cfg.ToolNaming = fileConfig.ToolNaming
			}
			if fileConfig.Shell != "" {
				cfg.Shell = fileConfig.Shell
			}
			if fileConfig.Lang != "" {
				cfg.Lang = fileConfig.Lang
			}
			if fileConfig.Widgets != "" {
				cfg.Widgets = fileConfig.Widgets
			}
			if len(fileConfig.SkillPaths) > 0 {
				cfg.SkillPaths = fileConfig.SkillPaths
			}
			if len(fileConfig.AllowedHosts) > 0 {
				cfg.AllowedHosts = fileConfig.AllowedHosts
			}
			if _, ok := raw["skillsEnabled"]; ok {
				cfg.SkillsEnabled = fileConfig.SkillsEnabled
			}
			if fileConfig.OAuth.OwnerToken != "" {
				cfg.OAuth.OwnerToken = fileConfig.OAuth.OwnerToken
			}
			if fileConfig.OAuth.AccessTokenTTLSeconds != 0 {
				cfg.OAuth.AccessTokenTTLSeconds = fileConfig.OAuth.AccessTokenTTLSeconds
			}
			if fileConfig.OAuth.RefreshTokenTTLSeconds != 0 {
				cfg.OAuth.RefreshTokenTTLSeconds = fileConfig.OAuth.RefreshTokenTTLSeconds
			}
			if len(fileConfig.OAuth.Scopes) > 0 {
				cfg.OAuth.Scopes = fileConfig.OAuth.Scopes
			}
			if len(fileConfig.OAuth.AllowedRedirectHosts) > 0 {
				cfg.OAuth.AllowedRedirectHosts = fileConfig.OAuth.AllowedRedirectHosts
			}
			if fileConfig.Logging.Level != "" {
				cfg.Logging.Level = fileConfig.Logging.Level
			}
			if fileConfig.Logging.Format != "" {
				cfg.Logging.Format = fileConfig.Logging.Format
			}
			if rawLogging, ok := raw["logging"]; ok {
				var loggingRaw map[string]json.RawMessage
				_ = json.Unmarshal(rawLogging, &loggingRaw)
				if _, ok := loggingRaw["requests"]; ok {
					cfg.Logging.Requests = fileConfig.Logging.Requests
				}
				if _, ok := loggingRaw["assets"]; ok {
					cfg.Logging.Assets = fileConfig.Logging.Assets
				}
				if _, ok := loggingRaw["toolCalls"]; ok {
					cfg.Logging.ToolCalls = fileConfig.Logging.ToolCalls
				}
				if _, ok := loggingRaw["shellCommands"]; ok {
					cfg.Logging.ShellCommands = fileConfig.Logging.ShellCommands
				}
				if _, ok := loggingRaw["trustProxy"]; ok {
					cfg.Logging.TrustProxy = fileConfig.Logging.TrustProxy
				}
			}
		}
	}

	// Check for owner password (new path first, old path as migration fallback)
	if cfg.OAuth.OwnerToken == "" {
		authFile := filepath.Join(cfg.ConfigDir, "auth.json")
		if _, err := os.Stat(authFile); os.IsNotExist(err) {
			oldAuthFile := filepath.Join(filepath.Dir(cfg.ConfigDir), ".devspace", "auth.json")
			if _, err := os.Stat(oldAuthFile); err == nil {
				authFile = oldAuthFile
			}
		}
		if data, err := os.ReadFile(authFile); err == nil {
			data = stripBOM(data)
			var auth struct {
				OwnerToken string `json:"ownerToken"`
			}
			if err := json.Unmarshal(data, &auth); err == nil && auth.OwnerToken != "" {
				cfg.OAuth.OwnerToken = auth.OwnerToken
			}
		}
	}

	return cfg
}

// ShellCommand returns the appropriate shell command for the current OS.
func (c *Config) ShellCommand() string {
	if v := os.Getenv("WEBCODER_SHELL"); v != "" {
		return v
	}
	if runtime.GOOS == "windows" {
		return "powershell.exe"
	}
	return "bash"
}

// ShellArgs returns shell arguments for command execution.
func (c *Config) ShellArgs(command string) []string {
	if runtime.GOOS == "windows" {
		return []string{"-NoProfile", "-NonInteractive", "-Command", command}
	}
	return []string{"-c", command}
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// stripBOM removes UTF-8 BOM (Byte Order Mark) from data.
// PowerShell's Set-Content -Encoding UTF8 adds a BOM that breaks json.Unmarshal.
func stripBOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	return data
}

// exeDir returns the directory containing the running executable.
// Falls back to current directory if detection fails.
func exeDir() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Dir(exe)
	}
	return "."
}

// detectSystemLang tries to detect the OS language.
// Returns a 2-letter code like "pl", "en", "de" or empty string.
func detectSystemLang() string {
	// Check LANG env (Linux/macOS)
	if lang := os.Getenv("LANG"); lang != "" {
		if len(lang) >= 2 {
			return strings.ToLower(lang[:2])
		}
	}
	if lang := os.Getenv("LC_ALL"); lang != "" {
		if len(lang) >= 2 {
			return strings.ToLower(lang[:2])
		}
	}

	// On Windows, try PowerShell
	if runtime.GOOS == "windows" {
		cmd := execCmd("powershell.exe", "-NoProfile", "-Command", "(Get-Culture).TwoLetterISOLanguageName")
		out, err := cmd.Output()
		if err == nil {
			code := strings.TrimSpace(string(out))
			if len(code) == 2 {
				return strings.ToLower(code)
			}
		}
	}

	return ""
}

// execCmd runs a command and returns the Cmd for reading output.
func execCmd(name string, args ...string) *execCmdResult {
	return &execCmdResult{name: name, args: args}
}

type execCmdResult struct {
	name string
	args []string
}

func (c *execCmdResult) Output() ([]byte, error) {
	cmd := exec.Command(c.name, c.args...)
	return cmd.Output()
}
