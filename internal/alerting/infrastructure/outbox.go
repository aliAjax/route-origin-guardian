package infrastructure

import (
	"context"
	"errors"
	"sync"

	"github.com/routeorigin/route-origin-guardian/internal/alerting/domain"
)

var ErrAlertIDRequired = errors.New("alert ID is required")

type Outbox struct {
	mu    sync.RWMutex
	items map[string]domain.Alert
}

func (o *Outbox) Enqueue(ctx context.Context, alert domain.Alert) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if alert.ID == "" {
		return ErrAlertIDRequired
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	o.items[alert.ID] = alert.Clone()
	return nil
}

func (o *Outbox) Get(id string) (domain.Alert, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	alert, ok := o.items[id]
	return alert.Clone(), ok
}
