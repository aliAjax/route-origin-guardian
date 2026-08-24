package application_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/routeorigin/route-origin-guardian/internal/refresh/adapter"
	"github.com/routeorigin/route-origin-guardian/internal/refresh/application"
	"github.com/routeorigin/route-origin-guardian/internal/refresh/domain"
	"github.com/routeorigin/route-origin-guardian/internal/refresh/infrastructure"
)

func TestRefreshLifecycleHonorsCancellation(t *testing.T) {
	t.Run("blocked source returns cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		source := infrastructure.NewSource(make(chan struct{}), domain.Snapshot{Provider: "rtr-a"})
		result := make(chan error, 1)
		go func() {
			_, err := source.Fetch(ctx)
			result <- err
		}()
		select {
		case err := <-result:
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("source error = %v, want context.Canceled", err)
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatal("source remained blocked after cancellation")
		}
	})

	t.Run("canceled commit leaves store unchanged", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		store := infrastructure.NewStore()
		err := store.Commit(ctx, domain.Snapshot{Provider: "rtr-a", Serial: 9})
		if !errors.Is(err, context.Canceled) || store.Count() != 0 {
			t.Fatalf("commit error = %v, snapshots = %d", err, store.Count())
		}
	})

	t.Run("coordinator does not detach caller context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		source := sourceFunc(func(got context.Context) (domain.Snapshot, error) {
			if got.Err() == nil {
				t.Fatal("source received a context without caller cancellation")
			}
			return domain.Snapshot{}, got.Err()
		})
		store := &recordingStore{}
		err := application.NewCoordinator(source, store).Refresh(ctx)
		if !errors.Is(err, context.Canceled) || store.commits != 0 {
			t.Fatalf("refresh error = %v, commits = %d", err, store.commits)
		}
	})

	t.Run("adapter recognizes wrapped context causes", func(t *testing.T) {
		canceled := adapter.Classify(fmt.Errorf("refresh stopped: %w", context.Canceled))
		deadline := adapter.Classify(fmt.Errorf("refresh stopped: %w", context.DeadlineExceeded))
		if canceled.Status != "canceled" || canceled.Retryable {
			t.Fatalf("canceled result = %+v", canceled)
		}
		if deadline.Status != "timed_out" || !deadline.Retryable {
			t.Fatalf("deadline result = %+v", deadline)
		}
	})
}

type sourceFunc func(context.Context) (domain.Snapshot, error)

func (f sourceFunc) Fetch(ctx context.Context) (domain.Snapshot, error) {
	return f(ctx)
}

type recordingStore struct {
	commits int
}

func (s *recordingStore) Commit(context.Context, domain.Snapshot) error {
	s.commits++
	return nil
}
