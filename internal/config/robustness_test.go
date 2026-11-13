package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigWarnsAboutMalformedFile(t *testing.T) {
	configDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(`{"port": 7676, "allowedRoots": [`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
	t.Setenv("DEVSPACE_LANG", "en")

	cfg := LoadConfig()
	if len(cfg.Warnings) != 1 || !strings.Contains(cfg.Warnings[0], "invalid JSON") {
		t.Fatalf("warnings = %#v", cfg.Warnings)
	}
	if cfg.Port != 7676 || len(cfg.AllowedRoots) != 0 {
		t.Fatalf("malformed file partially applied: %+v", cfg)
	}
}

func TestLoadConfigWithoutFileHasNoWarnings(t *testing.T) {
	t.Setenv("DEVSPACE_CONFIG_DIR", t.TempDir())
	t.Setenv("DEVSPACE_LANG", "en")
	t.Setenv("PORT", "")
	if cfg := LoadConfig(); len(cfg.Warnings) != 0 {
		t.Fatalf("warnings = %#v", cfg.Warnings)
	}
}

func TestLoadConfigWarnsAboutNonNumericPort(t *testing.T) {
	t.Setenv("DEVSPACE_CONFIG_DIR", t.TempDir())
	t.Setenv("DEVSPACE_LANG", "en")
	t.Setenv("PORT", "eighty")
	cfg := LoadConfig()
	if cfg.Port != 7676 || len(cfg.Warnings) != 1 {
		t.Fatalf("port=%d warnings=%#v", cfg.Port, cfg.Warnings)
	}
}

func TestConfigProblems(t *testing.T) {
	existing := t.TempDir()
	file := filepath.Join(existing, "file.txt")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := DefaultConfig()
	cfg.Port = 99999
	cfg.AllowedRoots = []string{existing, filepath.Join(existing, "missing"), file, "relative/dir"}
	problems := strings.Join(cfg.Problems(), "\n")
	for _, want := range []string{"port 99999", "missing does not exist", "not a directory", "not an absolute path"} {
		if !strings.Contains(problems, want) {
			t.Errorf("problems missing %q:\n%s", want, problems)
		}
	}
	if strings.Contains(problems, existing+" ") {
		t.Errorf("valid root reported:\n%s", problems)
	}
}

func TestWriteFileAtomicReplacesWithoutLeftovers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileAtomic(path, []byte(`{"new":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != `{"new":true}` {
		t.Fatalf("content = %q, %v", data, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("temporary files left behind: %v", entries)
	}
}

func TestWriteFileAtomicKeepsOldFileOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Renaming onto a directory fails after the data was written.
	target := filepath.Join(dir, "subdir")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "keep"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileAtomic(target, []byte("new"), 0o600); err == nil {
		t.Fatal("expected rename onto a non-empty directory to fail")
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 2 {
		t.Fatalf("temporary file left behind after failure: %v", entries)
	}
	if data, _ := os.ReadFile(path); string(data) != "old" {
		t.Fatalf("unrelated file changed: %q", data)
	}
}
