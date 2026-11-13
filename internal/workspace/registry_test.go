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
