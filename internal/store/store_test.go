package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCreatesDevSpaceDatabase(t *testing.T) {
	stateDir := t.TempDir()
	s, err := New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(stateDir, "devspace.db")); err != nil {
		t.Fatalf("expected devspace.db to be created: %v", err)
	}
}

func TestNewUsesLegacyDatabaseWhenPresent(t *testing.T) {
	stateDir := t.TempDir()
	s, err := New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.CreateSession(&WorkspaceSession{ID: "legacy", Root: "C:/project", Mode: "checkout"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	newPath := filepath.Join(stateDir, "devspace.db")
	legacyPath := filepath.Join(stateDir, "webcoder.db")
	if err := os.Rename(newPath, legacyPath); err != nil {
		t.Fatal(err)
	}

	s, err = New(stateDir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	session, err := s.GetSession("legacy")
	if err != nil {
		t.Fatalf("expected session from legacy database: %v", err)
	}
	if session.Root != "C:/project" {
		t.Fatalf("unexpected legacy session root: %q", session.Root)
	}
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected legacy database to be copied to devspace.db: %v", err)
	}
	if _, err := os.Stat(legacyPath); err != nil {
		t.Fatalf("expected legacy database to remain available for rollback: %v", err)
	}
}

func TestNewMigratesSiblingWebCoderStateDirectory(t *testing.T) {
	root := t.TempDir()
	legacyStateDir := filepath.Join(root, ".webcoder-state")
	if err := os.MkdirAll(legacyStateDir, 0700); err != nil {
		t.Fatal(err)
	}

	legacyStore, err := New(legacyStateDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := legacyStore.CreateSession(&WorkspaceSession{ID: "sibling", Root: "C:/portable", Mode: "checkout"}); err != nil {
		t.Fatal(err)
	}
	if err := legacyStore.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(legacyStateDir, "devspace.db"), filepath.Join(legacyStateDir, "webcoder.db")); err != nil {
		t.Fatal(err)
	}

	newStateDir := filepath.Join(root, ".devspace-state")
	s, err := New(newStateDir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if _, err := s.GetSession("sibling"); err != nil {
		t.Fatalf("expected session migrated from sibling legacy directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(newStateDir, "devspace.db")); err != nil {
		t.Fatalf("expected migrated devspace.db: %v", err)
	}
}
