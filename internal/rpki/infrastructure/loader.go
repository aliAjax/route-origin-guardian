package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/adapter"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/application"
	"sync"
	"time"
)

type Loader struct {
	source   adapter.Source
	service  *application.Service
	interval time.Duration
	mu       sync.Mutex
	lastErr  error
}

func (l *Loader) Load(ctx context.Context) error { return l.load(ctx) }

func NewLoader(src adapter.Source, svc *application.Service, interval time.Duration) *Loader {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Loader{source: src, service: svc, interval: interval}
}
func (l *Loader) Run(ctx context.Context) error {
	tick := time.NewTicker(l.interval)
	defer tick.Stop()
	for {
		if e := l.load(ctx); e != nil && !errors.Is(e, context.Canceled) {
			l.mu.Lock()
			l.lastErr = e
			l.mu.Unlock()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
		}
	}
}
func (l *Loader) load(ctx context.Context) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	snap, e := l.source.Fetch(ctx)
	if e != nil {
		return fmt.Errorf("fetch ROA: %w", e)
	}
	if e = l.service.Import(ctx, snap); e != nil {
		return fmt.Errorf("import ROA: %w", e)
	}
	return nil
}
func (l *Loader) LastError() error { l.mu.Lock(); defer l.mu.Unlock(); return l.lastErr }
