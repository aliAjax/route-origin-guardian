package application

import (
	"context"
	"errors"
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
func (m *Manager) Put(ctx context.Context, s *domain.Session) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s == nil {
		return errors.New("session is nil")
	}
	owned := *s
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[s.ID] = &owned
	return nil
}
func (m *Manager) Get(ctx context.Context, id string) (domain.Session, bool) {
	if ctx.Err() != nil {
		return domain.Session{}, false
	}
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
func (m *Manager) ExpireDue() []string {
	if m.clock == nil {
		return nil
	}
	return m.Expire(m.clock.Now())
}
func (m *Manager) Ready(now time.Time) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]string, 0)
	for id, s := range m.sessions {
		if s.State == domain.Failed && !now.IsZero() {
			ids = append(ids, id)
		}
	}
	return ids
}
func (m *Manager) Reopen(ctx context.Context, id string, now time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[id]
	if !ok {
		return errors.New("session not found")
	}
	return s.Open(now)
}
