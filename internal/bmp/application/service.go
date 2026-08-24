package application

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/rib/application"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
)

type Service struct{ rib *application.Service }

func New(r *application.Service) *Service { return &Service{rib: r} }
func (s *Service) Route(ctx context.Context, r domain.Route) (domain.Event, error) {
	return s.rib.Announce(ctx, r)
}
