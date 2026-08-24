package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/routeorigin/route-origin-guardian/internal/session/adapter"
	"github.com/routeorigin/route-origin-guardian/internal/session/application"
	"github.com/routeorigin/route-origin-guardian/internal/session/domain"
	"github.com/routeorigin/route-origin-guardian/internal/session/infrastructure"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func TestSessionLifecycleCoordinatesExpiryAndRetry(t *testing.T) {
	start := time.Unix(1_700_000_000, 0)
	keepalive := adapter.Keepalive{Interval: 10 * time.Second}
	if keepalive.Expired(start, start.Add(9*time.Second)) {
		t.Fatal("keepalive expired before its deadline")
	}
	if !keepalive.Expired(start, start.Add(10*time.Second)) {
		t.Fatal("keepalive remained active at its deadline")
	}

	policy := infrastructure.RetryPolicy{Base: time.Second, Max: 4 * time.Second}
	if got := policy.Delay(1); got != time.Second {
		t.Fatalf("first retry delay = %s, want 1s", got)
	}
	if got := policy.Delay(5); got != 4*time.Second {
		t.Fatalf("capped retry delay = %s, want 4s", got)
	}

	s := &domain.Session{ID: "peer-a", State: domain.Established, LastMessage: start, Hold: 10 * time.Second}
	m := application.NewManager(fixedClock{now: start.Add(10 * time.Second)})
	if err := m.Put(context.Background(), s); err != nil {
		t.Fatal(err)
	}
	if got := m.ExpireDue(); len(got) != 1 || got[0] != s.ID {
		t.Fatalf("expired sessions = %v, want [%s]", got, s.ID)
	}
	failed, ok := m.Get(context.Background(), s.ID)
	if !ok || failed.State != domain.Failed {
		t.Fatalf("session after expiry = %+v, found=%v", failed, ok)
	}
	if got := m.Ready(failed.RetryAt.Add(-time.Nanosecond)); len(got) != 0 {
		t.Fatalf("session became ready before retry deadline: %v", got)
	}
	if err := m.Reopen(context.Background(), s.ID, failed.RetryAt.Add(-time.Nanosecond)); !errors.Is(err, domain.ErrRetryNotReady) {
		t.Fatalf("early reopen error = %v", err)
	}
	if err := m.Reopen(context.Background(), s.ID, failed.RetryAt); err != nil {
		t.Fatalf("reopen at retry deadline: %v", err)
	}
}
