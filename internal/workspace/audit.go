package workspace

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	auditFileName = "outside-access.log"
	// auditMaxBytes rotates the log to outside-access.log.1 once exceeded.
	auditMaxBytes = 5 << 20
)

type auditEntry struct {
	Time   string `json:"time"`
	Origin string `json:"origin"`
	Tool   string `json:"tool"`
	Op     string `json:"op"`
	Input  string `json:"input"`
	Path   string `json:"path,omitempty"`
	Allow  bool   `json:"allow"`
	Reason string `json:"reason,omitempty"`
}

// auditLog appends one JSON line per outside-workspace access attempt. It
// records paths and decisions, never file contents.
type auditLog struct {
	mu   sync.Mutex
	path string
	now  func() time.Time
}

func newAuditLog(path string) *auditLog {
	return &auditLog{path: path, now: time.Now}
}

func (a *auditLog) record(entry auditEntry) {
	event := log.Info()
	if !entry.Allow {
		event = log.Warn()
	}
	event.Str("origin", entry.Origin).Str("tool", entry.Tool).Str("op", entry.Op).
		Str("path", entry.Path).Bool("allow", entry.Allow).Str("reason", entry.Reason).
		Msg("outside_workspace_access")
	if a == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	entry.Time = a.now().UTC().Format(time.RFC3339Nano)
	line, err := json.Marshal(entry)
	if err != nil {
		return
	}
	if info, err := os.Stat(a.path); err == nil && info.Size()+int64(len(line)) > auditMaxBytes {
		_ = os.Rename(a.path, a.path+".1")
	}
	file, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		log.Warn().Err(err).Msg("cannot write outside-access audit log")
		return
	}
	defer file.Close()
	_, _ = file.Write(append(line, '\n'))
}
