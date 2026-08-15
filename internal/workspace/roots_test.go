package workspace

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createWalkFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{
		filepath.Join(root, ".hermes", "AGENTS.md"),
		filepath.Join(root, "project", "AGENTS.md"),
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("test"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestUnixWorkspaceWalkSkipsHeavyDirectories(t *testing.T) {
	root := createWalkFixture(t)
	var visited []string
	err := walkWorkspaceUnix(root, func(path string, _ os.FileInfo) error {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		visited = append(visited, filepath.ToSlash(rel))
		return nil
	}, 100, time.Minute, time.Now)
	if err != nil {
		t.Fatal(err)
	}

	if len(visited) != 1 || visited[0] != "project/AGENTS.md" {
		t.Fatalf("Unix visited %#v, want only project/AGENTS.md", visited)
	}
}

func TestWindowsWorkspaceWalkKeepsOriginalDirectorySet(t *testing.T) {
	root := createWalkFixture(t)
	var foundHermes bool
	err := walkWorkspaceWindows(root, func(path string, _ os.FileInfo) error {
		if filepath.ToSlash(path) == filepath.ToSlash(filepath.Join(root, ".hermes", "AGENTS.md")) {
			foundHermes = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !foundHermes {
		t.Fatal("Windows traversal changed: .hermes should remain visible")
	}
}

func TestUnixWorkspaceWalkHonorsFileLimit(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 5; i++ {
		path := filepath.Join(root, string(rune('a'+i))+".txt")
		if err := os.WriteFile(path, []byte("test"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	visited := 0
	err := walkWorkspaceUnix(root, func(string, os.FileInfo) error {
		visited++
		return nil
	}, 2, time.Minute, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	if visited != 2 {
		t.Fatalf("visited %d files, want 2", visited)
	}
}

func TestUnixWorkspaceWalkHonorsTimeout(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}

	started := time.Unix(1, 0)
	clockCalls := 0
	now := func() time.Time {
		clockCalls++
		if clockCalls == 1 {
			return started
		}
		return started.Add(2 * time.Second)
	}

	visited := 0
	err := walkWorkspaceUnix(root, func(string, os.FileInfo) error {
		visited++
		return nil
	}, 100, time.Second, now)
	if err != nil {
		t.Fatal(err)
	}
	if visited != 0 {
		t.Fatalf("visited %d files after timeout, want 0", visited)
	}
}

func TestUnixWorkspaceWalkAllowsExplicitSkippedRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".config")
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}

	visited := 0
	err := walkWorkspaceUnix(root, func(string, os.FileInfo) error {
		visited++
		return nil
	}, 100, time.Minute, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	if visited != 1 {
		t.Fatalf("visited %d files, want 1", visited)
	}
}
