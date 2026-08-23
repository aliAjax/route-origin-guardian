package grpc

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/rib/application"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/domain"
)

type Server struct {
	routes *application.Service
	roas   *domain.Index
}

func NewServer(routes *application.Service, roas *domain.Index) *Server { return &Server{routes, roas} }
func (s *Server) Validate(_ context.Context, prefix string, asn uint32) domain.Result {
	return s.roas.Validate(prefix, asn)
}
func (s *Server) Subscribe(ctx context.Context, after uint64, limit int) ([]applicationEvent, error) {
	items, e := s.routes.Events(ctx, after, limit)
	out := make([]applicationEvent, len(items))
	for i, v := range items {
		out[i] = applicationEvent{ID: v.ID, Prefix: v.Route.Prefix, Operation: string(v.Route.Operation)}
	}
	return out, e
}

type applicationEvent struct {
	ID                uint64
	Prefix, Operation string
}
