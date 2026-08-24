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
	e.Route.ASPath = append([]uint32(nil), e.Route.ASPath...)
	e.Route.Communities = append([]string(nil), e.Route.Communities...)
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
	out := make([]domain.Event, limit)
	for i := range d.items[:limit] {
		out[i] = d.items[i]
		out[i].Route.ASPath = append([]uint32(nil), d.items[i].Route.ASPath...)
		out[i].Route.Communities = append([]string(nil), d.items[i].Route.Communities...)
	}
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
