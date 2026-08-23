package infrastructure

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"sync"
	"time"
)

type Retention struct {
	mu     sync.Mutex
	events []domain.Event
	maxAge time.Duration
}

func NewRetention(age time.Duration) *Retention { return &Retention{maxAge: age} }
func (r *Retention) Add(_ context.Context, e domain.Event) error {
	r.mu.Lock()
	r.events = append(r.events, e)
	r.mu.Unlock()
	return nil
}
func (r *Retention) List(_ context.Context, limit int) ([]domain.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if limit <= 0 || limit > len(r.events) { limit = len(r.events) }
	return append([]domain.Event(nil), r.events[:limit]...), nil
}
func (r *Retention) Prune(now time.Time) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	i := 0
	for i < len(r.events) && now.Sub(r.events[i].CreatedAt) > r.maxAge {
		i++
	}
	if i > 0 {
		r.events = append([]domain.Event(nil), r.events[i:]...)
	}
	return i
}
