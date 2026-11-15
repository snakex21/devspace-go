package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigReadsOutsideAccess(t *testing.T) {
	configDir := t.TempDir()
	data := `{"outsideAccess":{"roots":[{"path":"/srv/shared","access":"read"},{"path":"/srv/rw","access":"readwrite"}],"remote":false,"deny":["/srv/shared/private"],"audit":false}}`
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)

	cfg := LoadConfig()
	access := cfg.OutsideAccess
	if len(access.Roots) != 2 || access.Roots[1].Access != AccessReadWrite || access.RemoteEnabled() || len(access.Deny) != 1 || access.AuditEnabled() {
		t.Fatalf("outsideAccess not loaded: %+v", access)
	}
}

func TestOutsideAccessDefaults(t *testing.T) {
	configDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(`{"port":7676}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DEVSPACE_CONFIG_DIR", configDir)

	cfg := LoadConfig()
	// No roots, so nothing is reachable; once roots are added they also work
	// through the tunnel unless "remote": false is set.
	if len(cfg.OutsideAccess.Roots) != 0 || !cfg.OutsideAccess.RemoteEnabled() || !cfg.OutsideAccess.AuditEnabled() {
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

func TestOutsideAccessRemoteDefaultsOnWithRootsAndCanBeDisabled(t *testing.T) {
	for data, want := range map[string]bool{
		`{"outsideAccess":{"roots":[{"path":"/srv/a"}]}}`:                true,
		`{"outsideAccess":{"roots":[{"path":"/srv/a"}],"remote":true}}`:  true,
		`{"outsideAccess":{"roots":[{"path":"/srv/a"}],"remote":false}}`: false,
	} {
		configDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("DEVSPACE_CONFIG_DIR", configDir)
		if got := LoadConfig().OutsideAccess.RemoteEnabled(); got != want {
			t.Errorf("%s: RemoteEnabled = %v, want %v", data, got, want)
		}
	}
}
