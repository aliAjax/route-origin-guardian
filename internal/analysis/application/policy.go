package application

import (
	"github.com/routeorigin/route-origin-guardian/internal/analysis/domain"
	rib "github.com/routeorigin/route-origin-guardian/internal/rib/domain"
)

type Policy struct {
	PrivateASN   bool
	MaxPath      int
	DefaultRoute bool
}

func Evaluate(p Policy, r rib.Route) []domain.Finding {
	f := domain.Analyze(r)
	if p.MaxPath > 0 && len(r.ASPath) > p.MaxPath {
		f = append(f, domain.Finding{Kind: "path_length", RouteID: r.ID, Detail: "path exceeds configured maximum"})
	}
	if !p.PrivateASN {
		out := f[:0]
		for _, x := range f {
			if x.Kind != "private_asn" {
				out = append(out, x)
			}
		}
		f = out
	}
	if !p.DefaultRoute {
		out := f[:0]
		for _, x := range f {
			if x.Kind != "default_route" {
				out = append(out, x)
			}
		}
		f = out
	}
	return f
}
