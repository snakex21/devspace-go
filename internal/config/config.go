package config

import (
	"encoding/json"
	"fmt"
	"os"
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

// Config holds all Dev Space Go server configuration.
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
		StateDir:      filepath.Join(exeDir, ".devspace-state"),
		WorktreeRoot:  filepath.Join(exeDir, ".devspace", "worktrees"),
		AgentDir:      filepath.Join(exeDir, ".codex"),
		ConfigDir:     filepath.Join(exeDir, ".devspace"),
		ToolMode:      ToolModeFull,
		ToolNaming:    NamingShort,
		Shell:         "auto",
		Lang:          "auto",
		Widgets:       WidgetFull,
		SkillsEnabled: true,
		AllowedHosts:  []string{"*"},
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
	loadedLegacyConfig := false

	// Environment variable overrides
	if v := os.Getenv("HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Port)
	}
	if v := envValue("ALLOWED_ROOTS"); v != "" {
		cfg.AllowedRoots = splitAndTrim(v, ",")
	}
	if v := envValue("PUBLIC_BASE_URL"); v != "" {
		cfg.PublicBaseURL = v
	}
	if v := envValue("STATE_DIR"); v != "" {
		cfg.StateDir = v
	}
	if v := envValue("WORKTREE_ROOT"); v != "" {
		cfg.WorktreeRoot = v
	}
	if v := envValue("AGENT_DIR"); v != "" {
		cfg.AgentDir = v
	}
	if v := envValue("CONFIG_DIR"); v != "" {
		cfg.ConfigDir = v
	}
	if v := envValue("TOOL_MODE"); v != "" {
		cfg.ToolMode = ToolMode(v)
	}
	if v := envValue("TOOL_NAMING"); v != "" {
		cfg.ToolNaming = ToolNaming(v)
	}
	if v := envValue("SHELL"); v != "" {
		cfg.Shell = v
	}
	if v := envValue("LANG"); v != "" {
		cfg.Lang = v
	}
	if v := envValue("WIDGETS"); v != "" {
		cfg.Widgets = WidgetMode(v)
	}
	if v := envValue("SKILLS"); v == "0" {
		cfg.SkillsEnabled = false
	}
	if v := envValue("SKILL_PATHS"); v != "" {
		cfg.SkillPaths = splitAndTrim(v, ",")
	}
	if v := envValue("ALLOWED_HOSTS"); v != "" {
		cfg.AllowedHosts = splitAndTrim(v, ",")
	}

	// Logging config
	if v := envValue("LOG_LEVEL"); v != "" {
		cfg.Logging.Level = LogLevel(v)
	}
	if v := envValue("LOG_FORMAT"); v != "" {
		cfg.Logging.Format = LogFormat(v)
	}
	if v := envValue("LOG_REQUESTS"); v == "0" {
		cfg.Logging.Requests = false
	}
	if v := envValue("LOG_ASSETS"); v == "1" {
		cfg.Logging.Assets = true
	}
	if v := envValue("LOG_TOOL_CALLS"); v == "0" {
		cfg.Logging.ToolCalls = false
	}
	if v := envValue("LOG_SHELL_COMMANDS"); v == "1" {
		cfg.Logging.ShellCommands = true
	}
	if v := envValue("TRUST_PROXY"); v == "1" {
		cfg.Logging.TrustProxy = true
	}

	// Load from config file if exists (new path first, old path as migration fallback)
	configFile := filepath.Join(cfg.ConfigDir, "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Compatibility: use the legacy WebCoder config when no Dev Space Go
		// config exists yet. Saving from the CLI or GUI writes to .devspace.
		oldConfigFile := filepath.Join(filepath.Dir(cfg.ConfigDir), ".webcoder", "config.json")
		if _, err := os.Stat(oldConfigFile); err == nil {
			configFile = oldConfigFile
			loadedLegacyConfig = true
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

	if loadedLegacyConfig {
		rebaseLegacyPortablePaths(cfg)
	}
	resolvePortablePaths(cfg)

	// Resolve "auto" only after all sources have been merged. This avoids an
	// unnecessary OS lookup when the config file already specifies a language.
	if cfg.Lang == "auto" || cfg.Lang == "" {
		if detected := detectSystemLang(); detected != "" {
			cfg.Lang = detected
		} else {
			cfg.Lang = "en"
		}
	}

	return cfg
}

// PortablePath converts a path inside the application directory to a relative
// path suitable for config.json. External custom paths remain absolute.
func PortablePath(path string) string {
	if path == "" || !filepath.IsAbs(path) {
		return path
	}
	rel, err := filepath.Rel(exeDir(), path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return path
	}
	return rel
}

func resolvePortablePaths(cfg *Config) {
	base := exeDir()
	cfg.StateDir = resolvePortablePath(base, cfg.StateDir)
	cfg.WorktreeRoot = resolvePortablePath(base, cfg.WorktreeRoot)
	cfg.AgentDir = resolvePortablePath(base, cfg.AgentDir)
	cfg.ConfigDir = resolvePortablePath(base, cfg.ConfigDir)
}

func resolvePortablePath(base, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

// Legacy configs stored generated portable paths as absolute paths. Rebase
// those known defaults so moving the application folder also moves its data.
func rebaseLegacyPortablePaths(cfg *Config) {
	defaults := DefaultConfig()
	if pathEndsWith(cfg.StateDir, ".webcoder-state") || pathEndsWith(cfg.StateDir, ".devspace-state") {
		cfg.StateDir = defaults.StateDir
	}
	if pathEndsWith(cfg.WorktreeRoot, filepath.Join(".webcoder", "worktrees")) ||
		pathEndsWith(cfg.WorktreeRoot, filepath.Join(".devspace", "worktrees")) {
		cfg.WorktreeRoot = defaults.WorktreeRoot
	}
	if pathEndsWith(cfg.AgentDir, ".codex") {
		cfg.AgentDir = defaults.AgentDir
	}
}

func pathEndsWith(path, suffix string) bool {
	path = strings.ToLower(filepath.Clean(path))
	suffix = strings.ToLower(filepath.Clean(suffix))
	return path == suffix || strings.HasSuffix(path, string(filepath.Separator)+suffix)
}

// ShellCommand returns the appropriate shell command for the current OS.
func (c *Config) ShellCommand() string {
	if v := envValue("SHELL"); v != "" {
		return v
	}
	if runtime.GOOS == "windows" {
		return "powershell.exe"
	}
	return "bash"
}

// envValue returns a Dev Space Go environment variable, falling back to its
// legacy WebCoder equivalent for backwards compatibility.
func envValue(suffix string) string {
	if value := os.Getenv("DEVSPACE_" + suffix); value != "" {
		return value
	}
	return os.Getenv("WEBCODER_" + suffix)
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
	// Environment variables take precedence on every platform.
	if lang := os.Getenv("LC_ALL"); len(lang) >= 2 {
		return strings.ToLower(lang[:2])
	}
	if lang := os.Getenv("LANG"); lang != "" {
		if len(lang) >= 2 {
			return strings.ToLower(lang[:2])
		}
	}

	return detectOSSystemLang()
}
