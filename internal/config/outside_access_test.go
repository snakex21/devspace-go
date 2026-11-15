package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigReadsOutsideAccess(t *testing.T) {
	configDir := t.TempDir()
	data := `{"outsideAccess":{"roots":[{"path":"/srv/shared","access":"read"},{"path":"/srv/rw","access":"readwrite"}],"remote":true,"deny":["/srv/shared/private"],"audit":false}}`
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)

	cfg := LoadConfig()
	access := cfg.OutsideAccess
	if len(access.Roots) != 2 || access.Roots[1].Access != AccessReadWrite || !access.Remote || len(access.Deny) != 1 || access.AuditEnabled() {
		t.Fatalf("outsideAccess not loaded: %+v", access)
	}
}

func TestOutsideAccessDefaultsAreClosed(t *testing.T) {
	configDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(`{"port":7676}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)

	cfg := LoadConfig()
	if len(cfg.OutsideAccess.Roots) != 0 || cfg.OutsideAccess.Remote || !cfg.OutsideAccess.AuditEnabled() {
		t.Fatalf("unexpected defaults: %+v", cfg.OutsideAccess)
	}
}

func TestProblemsReportsInvalidExtraRootAccess(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutsideAccess.Roots = []ExtraRootConfig{{Path: "/x", Access: "write"}, {Path: "/y", Access: "READ"}}
	problems := strings.Join(cfg.Problems(), "\n")
	if !strings.Contains(problems, `"write"`) || strings.Contains(problems, `"READ"`) {
		t.Fatalf("problems: %s", problems)
	}
}
