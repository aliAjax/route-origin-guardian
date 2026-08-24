package application

import (
	"context"
	"fmt"

	"github.com/routeorigin/route-origin-guardian/internal/alerting/domain"
)

type Queue interface {
	Enqueue(context.Context, domain.Alert) error
}

type Sender interface {
	Send(context.Context, domain.Alert) error
}

type Dispatcher struct {
	queue  Queue
	sender Sender
}

func NewDispatcher(queue Queue, sender Sender) *Dispatcher {
	return &Dispatcher{queue: queue, sender: sender}
}

func (d *Dispatcher) Dispatch(ctx context.Context, alert domain.Alert) error {
	if err := d.queue.Enqueue(ctx, alert); err != nil {
		return fmt.Errorf("enqueue route-origin alert: %w", err)
	}
	if d.sender == nil {
		return nil
	}
	if err := d.sender.Send(ctx, alert.Clone()); err != nil {
		return fmt.Errorf("send route-origin alert: %w", err)
	}
	return nil
}
