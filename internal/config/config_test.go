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

	t.Setenv("WEBCODER_CONFIG_DIR", configDir)
	t.Setenv("WEBCODER_LANG", "")
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

	t.Setenv("WEBCODER_CONFIG_DIR", configDir)
	t.Setenv("WEBCODER_LANG", "")
	t.Setenv("LC_ALL", "pl_PL.UTF-8")
	t.Setenv("LANG", "")

	cfg := LoadConfig()
	if cfg.Lang != "de" {
		t.Fatalf("expected configured language de, got %q", cfg.Lang)
	}
}
