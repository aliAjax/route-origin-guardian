package infrastructure

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/rib/application"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"sort"
	"sync"
	"time"
)

type MemoryStore struct {
	mu     sync.RWMutex
	routes map[string]domain.Route
	events []domain.Event
	next   uint64
	max    int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{routes: make(map[string]domain.Route), max: 10000}
}
func (m *MemoryStore) Ready(context.Context) error { return nil }
func (m *MemoryStore) Upsert(_ context.Context, r domain.Route) (domain.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.routes[r.Key()] = r
	return m.eventLocked(r), nil
}
func (m *MemoryStore) Withdraw(_ context.Context, r domain.Route) (domain.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.routes, r.Key())
	return m.eventLocked(r), nil
}
func (m *MemoryStore) eventLocked(r domain.Route) domain.Event {
	m.next++
	e := domain.Event{ID: m.next, Route: r, CreatedAt: time.Now()}
	m.events = append(m.events, e)
	if len(m.events) > m.max {
		m.events = m.events[len(m.events)-m.max:]
	}
	return e
}
func (m *MemoryStore) List(_ context.Context, q application.Query) ([]domain.Route, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]domain.Route, 0)
	for _, r := range m.routes {
		if q.TenantID != "" && r.TenantID != q.TenantID {
			continue
		}
		if q.Prefix != "" && r.Prefix != q.Prefix {
			continue
		}
		if q.Peer != "" && r.PeerAddress != q.Peer {
			continue
		}
		if q.Family != "" && string(r.Family) != q.Family {
			continue
		}
		if q.ASN != 0 && r.OriginASN != q.ASN {
			continue
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key() < out[j].Key() })
	if len(out) > q.Limit {
		out = out[:q.Limit]
	}
	return out, nil
}
func (m *MemoryStore) Events(_ context.Context, after uint64, limit int) ([]domain.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]domain.Event, 0)
	for _, e := range m.events {
		if e.ID > after {
			out = append(out, e)
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}
