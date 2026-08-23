package application

import (
	"context"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"time"
)

type Store interface {
	Upsert(context.Context, domain.Route) (domain.Event, error)
	Withdraw(context.Context, domain.Route) (domain.Event, error)
	List(context.Context, Query) ([]domain.Route, error)
	Events(context.Context, uint64, int) ([]domain.Event, error)
	Ready(context.Context) error
}
type Query struct {
	TenantID, Prefix, Peer, Family string
	ASN                            uint32
	Limit                          int
}
type Service struct {
	store Store
	now   func() time.Time
}

func NewService(s Store) *Service { return &Service{store: s, now: time.Now} }
func (s *Service) Announce(ctx context.Context, r domain.Route) (domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return domain.Event{}, err
	}
	if err := validate(r); err != nil {
		return domain.Event{}, err
	}
	r.Operation = domain.Announce
	r.ReceivedAt = s.now()
	return s.store.Upsert(ctx, r)
}
func (s *Service) Replace(ctx context.Context, r domain.Route) (domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return domain.Event{}, err
	}
	if err := validate(r); err != nil {
		return domain.Event{}, err
	}
	r.Operation = domain.Replace
	r.ReceivedAt = s.now()
	return s.store.Upsert(ctx, r)
}
func (s *Service) Withdraw(ctx context.Context, r domain.Route) (domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return domain.Event{}, err
	}
	if r.Prefix == "" {
		return domain.Event{}, fmt.Errorf("withdraw prefix: %w", ErrInvalidRoute)
	}
	r.Operation = domain.Withdraw
	r.ReceivedAt = s.now()
	return s.store.Withdraw(ctx, r)
}
func (s *Service) List(ctx context.Context, q Query) ([]domain.Route, error) {
	if q.Limit <= 0 || q.Limit > 500 {
		q.Limit = 100
	}
	return s.store.List(ctx, q)
}
func (s *Service) Events(ctx context.Context, after uint64, limit int) ([]domain.Event, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	return s.store.Events(ctx, after, limit)
}
func (s *Service) Ready(ctx context.Context) error { return s.store.Ready(ctx) }

var ErrInvalidRoute = fmt.Errorf("invalid route")

func validate(r domain.Route) error {
	if r.Prefix == "" || r.PeerAddress == "" || r.ObserverID == "" || r.TenantID == "" {
		return fmt.Errorf("required route field: %w", ErrInvalidRoute)
	}
	if r.OriginASN == 0 {
		return fmt.Errorf("origin ASN: %w", ErrInvalidRoute)
	}
	return nil
}
