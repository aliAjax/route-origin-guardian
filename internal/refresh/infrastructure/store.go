package infrastructure

import (
	"context"
	"sync"

	"github.com/routeorigin/route-origin-guardian/internal/refresh/domain"
)

type Store struct {
	mu        sync.RWMutex
	snapshots []domain.Snapshot
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Commit(ctx context.Context, snapshot domain.Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots = append(s.snapshots, snapshot.Clone())
	return nil
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.snapshots)
}

func (s *Store) Latest() (domain.Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.snapshots) == 0 {
		return domain.Snapshot{}, false
	}
	return s.snapshots[len(s.snapshots)-1].Clone(), true
}
