package adapter

import (
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"strings"
)

type Filter struct {
	Tenant, Prefix, Peer string
	ASN                  uint32
}

func (f Filter) Match(r domain.Route) bool {
	if f.Tenant != "" && r.TenantID != f.Tenant {
		return false
	}
	if f.Prefix != "" && !strings.HasPrefix(r.Prefix, f.Prefix) {
		return false
	}
	if f.Peer != "" && r.PeerAddress != f.Peer {
		return false
	}
	if f.ASN != 0 && r.OriginASN != f.ASN {
		return false
	}
	return true
}
func NormalizePrefix(s string) string { return strings.TrimSpace(s) }
