package store

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// WorkspaceSession represents a persisted workspace session.
type WorkspaceSession struct {
	ID         string    `json:"id"`
	Root       string    `json:"root"`
	Mode       string    `json:"mode"`
	SourceRoot string    `json:"sourceRoot"`
	BaseRef    string    `json:"baseRef"`
	BaseSha    string    `json:"baseSha"`
	Managed    bool      `json:"managed"`
	CreatedAt  time.Time `json:"createdAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
}

// Store provides SQLite-based persistence for workspace sessions.
type Store struct {
	db *sql.DB
	mu sync.RWMutex
}

// New creates a new Store, initializing the database and schema.
func New(stateDir string) (*Store, error) {
	if err := os.MkdirAll(stateDir, 0700); err != nil {
		return nil, fmt.Errorf("create state dir: %w", err)
	}

	dbPath := databasePath(stateDir)
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Configure connection pool (SQLite works best with single writer)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return s, nil
}

func databasePath(stateDir string) string {
	dbPath := filepath.Join(stateDir, "devspace.db")
	if _, err := os.Stat(dbPath); err == nil {
		return dbPath
	}

	legacyPaths := []string{filepath.Join(stateDir, "webcoder.db")}
	if filepath.Base(filepath.Clean(stateDir)) == ".devspace-state" {
		legacyPaths = append(legacyPaths, filepath.Join(filepath.Dir(stateDir), ".webcoder-state", "webcoder.db"))
	}

	for _, legacyPath := range legacyPaths {
		if _, err := os.Stat(legacyPath); err != nil {
			continue
		}
		// A WAL may contain committed data not present in the main file. In that
		// case keep using the legacy database instead of risking a partial copy.
		if _, err := os.Stat(legacyPath + "-wal"); err == nil {
			return legacyPath
		}
		if err := copyDatabase(legacyPath, dbPath); err == nil {
			return dbPath
		}
		return legacyPath
	}

	return dbPath
}

func copyDatabase(sourcePath, destinationPath string) (err error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	temporary, err := os.CreateTemp(filepath.Dir(destinationPath), ".devspace-db-migration-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		temporary.Close()
		if err != nil {
			_ = os.Remove(temporaryPath)
		}
	}()

	if _, err = io.Copy(temporary, source); err != nil {
		return err
	}
	if err = temporary.Sync(); err != nil {
		return err
	}
	if err = temporary.Close(); err != nil {
		return err
	}
	if err = os.Chmod(temporaryPath, 0600); err != nil {
		return err
	}
	return os.Rename(temporaryPath, destinationPath)
}

// migrate creates the schema if it doesn't exist.
func (s *Store) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS workspace_sessions (
		id TEXT PRIMARY KEY,
		root TEXT NOT NULL,
		mode TEXT NOT NULL DEFAULT 'checkout',
		source_root TEXT DEFAULT '',
		base_ref TEXT DEFAULT '',
		base_sha TEXT DEFAULT '',
		managed INTEGER DEFAULT 0,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		last_used_at TEXT NOT NULL DEFAULT (datetime('now'))
	);

	CREATE INDEX IF NOT EXISTS idx_workspace_sessions_last_used
		ON workspace_sessions(last_used_at);
	`

	_, err := s.db.Exec(query)
	return err
}

// CreateSession inserts a new workspace session.
func (s *Store) CreateSession(session *WorkspaceSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(
		`INSERT INTO workspace_sessions (id, root, mode, source_root, base_ref, base_sha, managed, created_at, last_used_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		session.ID, session.Root, session.Mode, session.SourceRoot,
		session.BaseRef, session.BaseSha, boolToInt(session.Managed),
		now, now,
	)
	return err
}

// GetSession retrieves a workspace session by ID.
func (s *Store) GetSession(id string) (*WorkspaceSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session := &WorkspaceSession{}
	var managed int
	var createdAt, lastUsedAt string

	err := s.db.QueryRow(
		`SELECT id, root, mode, source_root, base_ref, base_sha, managed, created_at, last_used_at
		 FROM workspace_sessions WHERE id = ?`, id,
	).Scan(&session.ID, &session.Root, &session.Mode, &session.SourceRoot,
		&session.BaseRef, &session.BaseSha, &managed, &createdAt, &lastUsedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("unknown workspace session: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	session.Managed = managed != 0
	session.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	session.LastUsedAt, _ = time.Parse(time.RFC3339, lastUsedAt)

	return session, nil
}

// GetLatestSession retrieves the most recently used workspace session.
func (s *Store) GetLatestSession() (*WorkspaceSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session := &WorkspaceSession{}
	var managed int
	var createdAt, lastUsedAt string

	err := s.db.QueryRow(
		`SELECT id, root, mode, source_root, base_ref, base_sha, managed, created_at, last_used_at
		 FROM workspace_sessions
		 ORDER BY last_used_at DESC
		 LIMIT 1`,
	).Scan(&session.ID, &session.Root, &session.Mode, &session.SourceRoot,
		&session.BaseRef, &session.BaseSha, &managed, &createdAt, &lastUsedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no workspace sessions found")
	}
	if err != nil {
		return nil, fmt.Errorf("get latest session: %w", err)
	}

	session.Managed = managed != 0
	session.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	session.LastUsedAt, _ = time.Parse(time.RFC3339, lastUsedAt)

	return session, nil
}

// TouchSession updates the last_used_at timestamp.
func (s *Store) TouchSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`UPDATE workspace_sessions SET last_used_at = ? WHERE id = ?`,
		time.Now().UTC().Format(time.RFC3339), id,
	)
	return err
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
