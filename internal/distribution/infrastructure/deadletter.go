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

func (d *DeadLetter) Put(_ context.Context, e domain.Event) error {
	d.mu.Lock()
	d.items = append(d.items, e)
	d.mu.Unlock()
	return nil
}
func (d *DeadLetter) List(_ context.Context, limit int) []domain.Event {
	d.mu.Lock()
	defer d.mu.Unlock()
	if limit <= 0 || limit > len(d.items) {
		limit = len(d.items)
	}
	out := append([]domain.Event(nil), d.items[:limit]...)
	return out
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
