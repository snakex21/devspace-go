package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/store"
)

func TestRegistryConcurrentOpenAndLookup(t *testing.T) {
	root := t.TempDir()
	s, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}
	cfg.AgentDir = filepath.Join(t.TempDir(), "agents")
	registry := NewRegistry(cfg, s)

	var wg sync.WaitGroup
	errs := make(chan error, 64)
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				dir := filepath.Join(root, fmt.Sprintf("p%d", worker))
				ctx, err := registry.OpenWorkspace(dir, ModeCheckout, "")
				if err != nil {
					errs <- err
					return
				}
				for _, id := range []string{ctx.Workspace.ID, "default", "latest"} {
					if _, err := registry.GetWorkspace(id); err != nil {
						errs <- err
						return
					}
				}
			}
		}(worker)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
}

func TestRegistryRestoresWorkspaceAfterRestart(t *testing.T) {
	root := t.TempDir()
	stateDir := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}

	s, err := store.New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	opened, err := NewRegistry(cfg, s).OpenWorkspace(root, ModeCheckout, "")
	if err != nil {
		t.Fatal(err)
	}
	s.Close()

	s, err = store.New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ws, err := NewRegistry(cfg, s).GetWorkspace(opened.Workspace.ID)
	if err != nil || ws.Root != filepath.Clean(root) {
		t.Fatalf("restored %v, %v", ws, err)
	}
}

func TestRegistryRejectsStoredRootNoLongerAllowed(t *testing.T) {
	root := t.TempDir()
	stateDir := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}
	s, err := store.New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	opened, err := NewRegistry(cfg, s).OpenWorkspace(root, ModeCheckout, "")
	if err != nil {
		t.Fatal(err)
	}
	other := t.TempDir()
	cfg2 := config.DefaultConfig()
	cfg2.AllowedRoots = []string{other}
	if _, err := NewRegistry(cfg2, s).GetWorkspace(opened.Workspace.ID); err == nil {
		t.Fatal("workspace outside the new allowed roots was restored")
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultWorkspaceFollowsLastUsedAcrossRestart(t *testing.T) {
	parent := t.TempDir()
	stateDir := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{parent}
	cfg.AgentDir = filepath.Join(t.TempDir(), "agents")

	s, err := store.New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	registry := NewRegistry(cfg, s)
	a, err := registry.OpenWorkspace(filepath.Join(parent, "a"), ModeCheckout, "")
	if err != nil {
		t.Fatal(err)
	}
	b, err := registry.OpenWorkspace(filepath.Join(parent, "b"), ModeCheckout, "")
	if err != nil {
		t.Fatal(err)
	}
	if ws, _ := registry.GetWorkspace("default"); ws.ID != b.Workspace.ID {
		t.Fatalf("default = %s, want most recently opened b", ws.ID)
	}
	// Switching back to a must be reflected immediately and persisted, even
	// though last_used_at writes are throttled for repeated use.
	for i := 0; i < 5; i++ {
		if _, err := registry.GetWorkspace(a.Workspace.ID); err != nil {
			t.Fatal(err)
		}
	}
	if ws, _ := registry.GetWorkspace("latest"); ws.ID != a.Workspace.ID {
		t.Fatalf("latest = %s, want a", ws.ID)
	}
	s.Close()

	// The switch back to a was persisted: after a restart, a is at least as
	// recent as b (last_used_at has one-second resolution).
	s, err = store.New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if mustSession(t, s, a.Workspace.ID).LastUsedAt.Before(mustSession(t, s, b.Workspace.ID).LastUsedAt) {
		t.Fatal("switching back to a was not persisted")
	}
}

func mustSession(t *testing.T, s *store.Store, id string) *store.WorkspaceSession {
	t.Helper()
	session, err := s.GetSession(id)
	if err != nil {
		t.Fatal(err)
	}
	return session
}

func TestDefaultWorkspaceWithoutStoreReusesWorkspace(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{t.TempDir()}
	cfg.AgentDir = filepath.Join(t.TempDir(), "agents")
	registry := NewRegistry(cfg, nil)
	first, err := registry.GetWorkspace("default")
	if err != nil {
		t.Fatal(err)
	}
	second, err := registry.GetWorkspace("default")
	if err != nil || second.ID != first.ID {
		t.Fatalf("default created a new workspace on every call: %v %v", first.ID, second)
	}
}
