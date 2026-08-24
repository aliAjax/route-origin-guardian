package infrastructure

import (
	"context"
	"errors"

	"github.com/routeorigin/route-origin-guardian/internal/alerting/domain"
)

var ErrSenderUnavailable = errors.New("alert sender is unavailable")

type Sender interface {
	Send(context.Context, domain.Alert) error
}

type WebhookSender struct {
	endpoint string
}

func NewSender(enabled bool, endpoint string) Sender {
	if !enabled {
		var sender *WebhookSender
		return sender
	}
	return &WebhookSender{endpoint: endpoint}
}

func (s *WebhookSender) Send(ctx context.Context, _ domain.Alert) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s == nil || s.endpoint == "" {
		return ErrSenderUnavailable
	}
	return nil
}
