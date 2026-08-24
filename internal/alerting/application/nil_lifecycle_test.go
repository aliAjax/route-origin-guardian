package application_test

import (
	"context"
	"testing"

	"github.com/routeorigin/route-origin-guardian/internal/alerting/adapter"
	"github.com/routeorigin/route-origin-guardian/internal/alerting/application"
	"github.com/routeorigin/route-origin-guardian/internal/alerting/domain"
	"github.com/routeorigin/route-origin-guardian/internal/alerting/infrastructure"
)

func TestAlertDeliveryNilBoundaries(t *testing.T) {
	t.Run("alert without route metadata remains cloneable", func(t *testing.T) {
		alert := domain.Alert{ID: "alert-48", Labels: map[string]string{"origin": "invalid"}}
		if panicked(func() { _ = alert.Clone() }) {
			t.Fatal("cloning an alert without route metadata panicked")
		}
	})

	t.Run("zero value outbox accepts its first alert", func(t *testing.T) {
		var outbox infrastructure.Outbox
		alert := domain.Alert{ID: "alert-48"}
		if panicked(func() { _ = outbox.Enqueue(context.Background(), alert) }) {
			t.Fatal("zero value outbox panicked on first enqueue")
		}
		if _, ok := outbox.Get(alert.ID); !ok {
			t.Fatal("first alert was not retained")
		}
	})

	t.Run("disabled sender is a nil interface", func(t *testing.T) {
		if sender := infrastructure.NewSender(false, ""); sender != nil {
			t.Fatalf("disabled sender = %#v, want nil", sender)
		}
	})

	t.Run("dispatcher skips a typed nil sender", func(t *testing.T) {
		var sender *panicSender
		var outbox infrastructure.Outbox
		dispatcher := application.NewDispatcher(&outbox, sender)
		if panicked(func() {
			if err := dispatcher.Dispatch(context.Background(), domain.Alert{ID: "alert-48"}); err != nil {
				t.Fatalf("dispatch error = %v", err)
			}
		}) {
			t.Fatal("dispatcher invoked a typed nil sender")
		}
	})

	t.Run("typed nil delivery error means success", func(t *testing.T) {
		var deliveryErr *adapter.DeliveryError
		var err error = deliveryErr
		result := adapter.Classify(err)
		if result.Status != "delivered" || result.Retryable {
			t.Fatalf("typed nil result = %+v", result)
		}
	})
}

type panicSender struct{}

func (s *panicSender) Send(context.Context, domain.Alert) error {
	if s == nil {
		panic("typed nil sender invoked")
	}
	return nil
}

func panicked(fn func()) (didPanic bool) {
	defer func() {
		didPanic = recover() != nil
	}()
	fn()
	return false
}
