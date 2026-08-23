package domain

type Tenant struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
type Scope struct {
	TenantID  string
	Routers   map[string]bool
	Observers map[string]bool
}

func (s Scope) Allows(router, observer string) bool {
	return (len(s.Routers) == 0 || s.Routers[router]) && (len(s.Observers) == 0 || s.Observers[observer])
}
