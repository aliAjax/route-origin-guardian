package infrastructure

import (
	"context"

	"github.com/routeorigin/route-origin-guardian/internal/refresh/domain"
)

// Source models a refresh endpoint whose next snapshot may not be ready yet.
type Source struct {
	ready    <-chan struct{}
	snapshot domain.Snapshot
}

func NewSource(ready <-chan struct{}, snapshot domain.Snapshot) *Source {
	return &Source{ready: ready, snapshot: snapshot.Clone()}
}

func (s *Source) Fetch(ctx context.Context) (domain.Snapshot, error) {
	select {
	case <-ctx.Done():
		return domain.Snapshot{}, ctx.Err()
	case <-s.ready:
		return s.snapshot.Clone(), nil
	}
}
