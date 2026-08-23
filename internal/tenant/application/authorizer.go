package application

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/tenant/domain"
	"sync"
)

type Authorizer struct {
	mu     sync.RWMutex
	keys   map[string]string
	scopes map[string]domain.Scope
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{keys: make(map[string]string), scopes: make(map[string]domain.Scope)}
}
func (a *Authorizer) CreateKey(tenant, plain string) error {
	if tenant == "" || plain == "" {
		return fmt.Errorf("tenant and key required")
	}
	h := sha256.Sum256([]byte(plain))
	a.mu.Lock()
	a.keys[tenant] = hex.EncodeToString(h[:])
	a.mu.Unlock()
	return nil
}
func (a *Authorizer) Authenticate(tenant, plain string) bool {
	h := sha256.Sum256([]byte(plain))
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.keys[tenant] == hex.EncodeToString(h[:])
}
func (a *Authorizer) SetScope(s domain.Scope) { a.mu.Lock(); a.scopes[s.TenantID] = s; a.mu.Unlock() }
func (a *Authorizer) Allows(tenant, router, observer string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	s, ok := a.scopes[tenant]
	return ok && s.Allows(router, observer)
}
