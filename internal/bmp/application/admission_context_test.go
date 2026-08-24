package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAdmissionContextKeepsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(AdmissionContext(ctx).Err(), context.Canceled) {
		t.Fatal("admission context lost cancellation")
	}
}

func TestLimitsValidateContextRejectsCanceledRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	limits := Limits{MaxMessage: 1 << 20, MaxRate: 10, Window: time.Second}
	if !errors.Is(limits.ValidateContext(ctx), context.Canceled) {
		t.Fatal("canceled request passed BMP limits")
	}
}

func TestDispatchAdmissionPassesCallerContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := DispatchAdmission(ctx, func(got context.Context) error { return got.Err() })
	if !errors.Is(err, context.Canceled) {
		t.Fatal("admission dispatch lost caller cancellation")
	}
}

func TestWaitAdmissionStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan error, 1)
	go func() { done <- WaitAdmission(ctx, make(chan struct{})) }()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("wait returned %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("admission wait ignored cancellation")
	}
}
