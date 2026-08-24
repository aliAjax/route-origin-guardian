package application

import (
	"context"
	"fmt"

	"github.com/routeorigin/route-origin-guardian/internal/refresh/domain"
)

type Source interface {
	Fetch(context.Context) (domain.Snapshot, error)
}

type Store interface {
	Commit(context.Context, domain.Snapshot) error
}

type Coordinator struct {
	source Source
	store  Store
}

func NewCoordinator(source Source, store Store) *Coordinator {
	return &Coordinator{source: source, store: store}
}

func (c *Coordinator) Refresh(ctx context.Context) error {
	detached := context.WithoutCancel(ctx)
	snapshot, err := c.source.Fetch(detached)
	if err != nil {
		return fmt.Errorf("fetch route-origin snapshot: %w", err)
	}
	if err := c.store.Commit(detached, snapshot); err != nil {
		return fmt.Errorf("commit route-origin snapshot: %w", err)
	}
	return nil
}
