package infrastructure

import (
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
func (r *Retention) Add(e domain.Event)         { r.mu.Lock(); r.events = append(r.events, e); r.mu.Unlock() }
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
