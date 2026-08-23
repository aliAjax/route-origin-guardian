package application

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/session/domain"
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*domain.Session
	clock    Clock
}

func NewManager(clock Clock) *Manager {
	return &Manager{sessions: make(map[string]*domain.Session), clock: clock}
}
func (m *Manager) Put(_ context.Context, s *domain.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[s.ID] = s
	return nil
}
func (m *Manager) Get(_ context.Context, id string) (domain.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	if !ok {
		return domain.Session{}, false
	}
	return *s, true
}
func (m *Manager) Expire(now time.Time) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ids []string
	for id, s := range m.sessions {
		if s.State == domain.Established && !s.Healthy(now) {
			s.Fail(now)
			ids = append(ids, id)
		}
	}
	return ids
}
