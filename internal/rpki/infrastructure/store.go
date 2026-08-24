package infrastructure

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/application"
	"sync"
)

type SnapshotStore struct {
	mu    sync.RWMutex
	items []application.Snapshot
}

func (s *SnapshotStore) Save(_ context.Context, v application.Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, v)
	if len(s.items) > 32 {
		s.items = s.items[len(s.items)-32:]
	}
	return nil
}
func (s *SnapshotStore) Latest(_ context.Context) (application.Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return application.Snapshot{}, false
	}
	return s.items[len(s.items)-1], true
}
