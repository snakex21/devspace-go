package store

import (
	"database/sql"
	"fmt"
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

	dbPath := filepath.Join(stateDir, "webcoder.db")
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
