package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigResolvesAutoLanguageAfterFileMerge(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, []byte(`{"lang":"auto"}`), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
	t.Setenv("DEVSPACE_LANG", "")
	t.Setenv("LC_ALL", "de_DE.UTF-8")
	t.Setenv("LANG", "")

	cfg := LoadConfig()
	if cfg.Lang != "de" {
		t.Fatalf("expected auto language to resolve to de, got %q", cfg.Lang)
	}
}

func TestLoadConfigKeepsExplicitFileLanguage(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, []byte(`{"lang":"de"}`), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
	t.Setenv("DEVSPACE_LANG", "")
	t.Setenv("LC_ALL", "pl_PL.UTF-8")
	t.Setenv("LANG", "")

	cfg := LoadConfig()
	if cfg.Lang != "de" {
		t.Fatalf("expected configured language de, got %q", cfg.Lang)
	}
}

func TestDefaultConfigUsesDevSpacePortableDirectories(t *testing.T) {
	cfg := DefaultConfig()

	if filepath.Base(cfg.ConfigDir) != ".devspace" {
		t.Fatalf("expected .devspace config directory, got %q", cfg.ConfigDir)
	}
	if filepath.Base(cfg.StateDir) != ".devspace-state" {
		t.Fatalf("expected .devspace-state directory, got %q", cfg.StateDir)
	}
	if filepath.Base(filepath.Dir(cfg.WorktreeRoot)) != ".devspace" {
		t.Fatalf("expected worktrees under .devspace, got %q", cfg.WorktreeRoot)
	}
}

func TestLoadConfigPrefersDevSpaceEnvironmentVariables(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
	t.Setenv("DEVSPACE_ALLOWED_ROOTS", "C:/new-brand")
	t.Setenv("WEBCODER_ALLOWED_ROOTS", "C:/legacy-brand")

	cfg := LoadConfig()
	if len(cfg.AllowedRoots) != 1 || cfg.AllowedRoots[0] != "C:/new-brand" {
		t.Fatalf("expected DEVSPACE_ALLOWED_ROOTS to win, got %#v", cfg.AllowedRoots)
	}
}

func TestLoadConfigSupportsLegacyEnvironmentVariables(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
	t.Setenv("DEVSPACE_ALLOWED_ROOTS", "")
	t.Setenv("WEBCODER_ALLOWED_ROOTS", "C:/legacy-brand")

	cfg := LoadConfig()
	if len(cfg.AllowedRoots) != 1 || cfg.AllowedRoots[0] != "C:/legacy-brand" {
		t.Fatalf("expected legacy environment fallback, got %#v", cfg.AllowedRoots)
	}
}

func TestLoadConfigFallsBackToLegacyConfigFile(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, ".devspace")
	legacyDir := filepath.Join(root, ".webcoder")
	if err := os.MkdirAll(legacyDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacyDir, "config.json"), []byte(`{"lang":"fr"}`), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
	t.Setenv("WEBCODER_CONFIG_DIR", "")
	t.Setenv("DEVSPACE_LANG", "")
	t.Setenv("WEBCODER_LANG", "")

	cfg := LoadConfig()
	if cfg.Lang != "fr" {
		t.Fatalf("expected language from legacy config, got %q", cfg.Lang)
	}
	if cfg.ConfigDir != configDir {
		t.Fatalf("expected future saves to use %q, got %q", configDir, cfg.ConfigDir)
	}
}

func TestLoadConfigRebasesLegacyPortablePaths(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, ".devspace")
	legacyDir := filepath.Join(root, ".webcoder")
	if err := os.MkdirAll(legacyDir, 0700); err != nil {
		t.Fatal(err)
	}
	legacyConfig := `{
		"stateDir":"C:/previous-install/.webcoder-state",
		"worktreeRoot":"C:/previous-install/.webcoder/worktrees",
		"agentDir":"C:/previous-install/.codex"
	}`
	if err := os.WriteFile(filepath.Join(legacyDir, "config.json"), []byte(legacyConfig), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
	t.Setenv("WEBCODER_CONFIG_DIR", "")

	cfg := LoadConfig()
	defaults := DefaultConfig()
	if cfg.StateDir != defaults.StateDir || cfg.WorktreeRoot != defaults.WorktreeRoot || cfg.AgentDir != defaults.AgentDir {
		t.Fatalf("expected legacy paths to be rebased, got state=%q worktrees=%q agent=%q", cfg.StateDir, cfg.WorktreeRoot, cfg.AgentDir)
	}
}

func TestPortablePathKeepsApplicationDataRelative(t *testing.T) {
	cfg := DefaultConfig()
	if got := PortablePath(cfg.StateDir); got != ".devspace-state" {
		t.Fatalf("expected relative state path, got %q", got)
	}
	if got := PortablePath(cfg.WorktreeRoot); got != filepath.Join(".devspace", "worktrees") {
		t.Fatalf("expected relative worktree path, got %q", got)
	}

	external := filepath.Join(t.TempDir(), "state")
	if got := PortablePath(external); got != external {
		t.Fatalf("expected external path to remain absolute, got %q", got)
	}
}
