package infrastructure

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/tenant/domain"
	"sync"
)

type Repository struct {
	mu    sync.RWMutex
	items map[string]domain.Tenant
}

func NewRepository() *Repository { return &Repository{items: make(map[string]domain.Tenant)} }
func (r *Repository) Put(_ context.Context, t domain.Tenant) error {
	r.mu.Lock()
	r.items[t.ID] = t
	r.mu.Unlock()
	return nil
}
func (r *Repository) Get(_ context.Context, id string) (domain.Tenant, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	return t, ok
}
func (r *Repository) List(_ context.Context) []domain.Tenant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Tenant, 0, len(r.items))
	for _, t := range r.items {
		out = append(out, t)
	}
	return out
}
