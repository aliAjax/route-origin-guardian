package infrastructure

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"sync"
)

type DeadLetter struct {
	mu    sync.Mutex
	items []domain.Event
}

func (d *DeadLetter) Put(ctx context.Context, e domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = append(d.items, e)
	return nil
}
func (d *DeadLetter) List(ctx context.Context, limit int) []domain.Event {
	if ctx.Err() != nil {
		return nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if limit <= 0 || limit > len(d.items) {
		limit = len(d.items)
	}
	return append([]domain.Event(nil), d.items[:limit]...)
}
func (d *DeadLetter) Remove(id uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, e := range d.items {
		if e.ID == id {
			d.items = append(d.items[:i], d.items[i+1:]...)
			return
		}
	}
}
