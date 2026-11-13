package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestStoreConfiguresBusyTimeout(t *testing.T) {
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var timeout int
	if err := s.db.QueryRow("PRAGMA busy_timeout").Scan(&timeout); err != nil {
		t.Fatal(err)
	}
	if timeout < 1000 {
		t.Fatalf("busy_timeout = %d ms; concurrent writers fail immediately with SQLITE_BUSY", timeout)
	}
}

// Two Store instances on one state directory model two server processes, or
// a restarted server while the old one is still shutting down.
func TestStoreConcurrentWritersAcrossConnections(t *testing.T) {
	dir := t.TempDir()
	first, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()

	var wg sync.WaitGroup
	errs := make(chan error, 400)
	for worker := 0; worker < 8; worker++ {
		s := first
		if worker%2 == 1 {
			s = second
		}
		wg.Add(1)
		go func(worker int, s *Store) {
			defer wg.Done()
			for i := 0; i < 25; i++ {
				id := fmt.Sprintf("ws_%d_%d", worker, i)
				if err := s.CreateSession(&WorkspaceSession{ID: id, Root: "/tmp/x", Mode: "checkout"}); err != nil {
					errs <- err
					return
				}
				if err := s.TouchSession(id); err != nil {
					errs <- err
					return
				}
				if _, err := s.GetSession(id); err != nil {
					errs <- err
					return
				}
				if _, err := s.GetLatestSession(); err != nil {
					errs <- err
					return
				}
			}
		}(worker, s)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
}
