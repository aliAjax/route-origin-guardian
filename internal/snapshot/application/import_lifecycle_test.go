package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/routeorigin/route-origin-guardian/internal/snapshot/adapter"
	"github.com/routeorigin/route-origin-guardian/internal/snapshot/application"
	"github.com/routeorigin/route-origin-guardian/internal/snapshot/domain"
	"github.com/routeorigin/route-origin-guardian/internal/snapshot/infrastructure"
)

func TestSnapshotImportLifecycle(t *testing.T) {
	t.Run("batch values do not retain caller storage", func(t *testing.T) {
		prefixes := []string{"203.0.113.0/24"}
		batch, err := domain.NewBatch("epoch-7", prefixes)
		if err != nil {
			t.Fatal(err)
		}
		prefixes[0] = "198.51.100.0/24"
		if batch.Prefixes[0] != "203.0.113.0/24" {
			t.Fatalf("batch retained caller prefixes: %v", batch.Prefixes)
		}
		clone := batch.Clone()
		clone.Prefixes[0] = "192.0.2.0/24"
		if batch.Prefixes[0] != "203.0.113.0/24" {
			t.Fatalf("batch clone shares prefixes: %v", batch.Prefixes)
		}
	})

	t.Run("subscription finalization is idempotent", func(t *testing.T) {
		batch, _ := domain.NewBatch("epoch-8", []string{"203.0.113.0/24"})
		subscription, err := infrastructure.NewSource([]domain.Batch{batch}, nil).Open(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if err := subscription.Close(); err != nil {
			t.Fatal(err)
		}
		if err := subscription.Close(); err != nil {
			t.Fatal(err)
		}
		if subscription.CloseCount() != 1 {
			t.Fatalf("close count = %d", subscription.CloseCount())
		}
	})

	t.Run("rollback discards staged batches", func(t *testing.T) {
		batch, _ := domain.NewBatch("epoch-9", []string{"203.0.113.0/24"})
		transaction := infrastructure.NewStore().Begin()
		if err := transaction.Stage(context.Background(), batch); err != nil {
			t.Fatal(err)
		}
		if err := transaction.Rollback(); err != nil {
			t.Fatal(err)
		}
		if transaction.StagedCount() != 0 {
			t.Fatalf("rollback retained %d batches", transaction.StagedCount())
		}
		if err := transaction.Stage(context.Background(), batch); !errors.Is(err, infrastructure.ErrTransactionClosed) {
			t.Fatalf("stage after rollback error = %v", err)
		}
	})

	t.Run("failed import keeps storage empty", func(t *testing.T) {
		valid, _ := domain.NewBatch("epoch-10", []string{"203.0.113.0/24"})
		invalid := domain.Batch{ID: "epoch-11"}
		store := infrastructure.NewStore()
		importer := application.Importer{Source: infrastructure.NewSource([]domain.Batch{valid, invalid}, nil), Store: store}
		err := importer.Sync(context.Background())
		if !errors.Is(err, domain.ErrInvalidBatch) {
			t.Fatalf("import error = %v", err)
		}
		if got := store.Snapshot(); len(got) != 0 {
			t.Fatalf("failed import committed %d batches", len(got))
		}
	})

	t.Run("wrapped lifecycle causes keep their identity", func(t *testing.T) {
		if got := adapter.ClassifyImport(errors.Join(errors.New("sync failed"), context.Canceled)); got.Status != "canceled" {
			t.Fatalf("canceled result = %+v", got)
		}
		if got := adapter.ClassifyImport(errors.Join(errors.New("sync failed"), domain.ErrInvalidBatch)); got.Status != "invalid" {
			t.Fatalf("invalid result = %+v", got)
		}
	})
}
