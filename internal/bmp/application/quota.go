package application

import (
	"context"
	"time"
)

type Quota struct {
	limits Limits
	used   int
	start  time.Time
}

func NewQuota(l Limits, now time.Time) *Quota { return &Quota{limits: l, start: now} }

func (q *Quota) Reserve(ctx context.Context, size int, now time.Time) (*Lease, error) {
	q.advanceWindow(now)
	if !q.limits.Allow(size, q.used) {
		return nil, ErrQuotaExceeded
	}
	q.used += size
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &Lease{quota: q, size: size}, nil
}

var ErrQuotaExceeded = errorString("BMP quota exceeded")

type errorString string

func (e errorString) Error() string { return string(e) }
