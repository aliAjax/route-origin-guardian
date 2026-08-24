package application

import (
	"context"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
)

type EventSource interface {
	Events(context.Context, uint64, int) ([]domain.Event, error)
}
type Cursor struct{ source EventSource }

func NewCursor(s EventSource) *Cursor { return &Cursor{source: s} }
func (c *Cursor) Read(ctx context.Context, after uint64, limit int) ([]domain.Event, uint64, error) {
	if err := ctx.Err(); err != nil {
		return nil, after, err
	}
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	items, e := c.source.Events(ctx, after, limit)
	if e != nil {
		return nil, after, fmt.Errorf("read event cursor: %w", e)
	}
	next := after
	if len(items) > 0 {
		next = items[len(items)-1].ID
	}
	return items, next, nil
}
