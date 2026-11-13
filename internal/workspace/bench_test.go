package workspace

import (
	"path/filepath"
	"testing"

	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/store"
)

func BenchmarkResolveInRoot(b *testing.B) {
	root := b.TempDir()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := ResolveInRoot(root, filepath.Join("src", "pkg", "file.go")); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetWorkspaceWithStore(b *testing.B) {
	root := b.TempDir()
	s, err := store.New(b.TempDir())
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()
	cfg := config.DefaultConfig()
	cfg.AllowedRoots = []string{root}
	registry := NewRegistry(cfg, s)
	ws, err := registry.OpenWorkspace(root, ModeCheckout, "")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := registry.GetWorkspace(ws.Workspace.ID); err != nil {
			b.Fatal(err)
		}
	}
}
