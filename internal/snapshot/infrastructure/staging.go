package infrastructure

import (
	"context"
	"errors"
	"sync"

	"github.com/routeorigin/route-origin-guardian/internal/snapshot/domain"
)

var ErrTransactionClosed = errors.New("snapshot transaction closed")

type Store struct {
	mu      sync.RWMutex
	batches []domain.Batch
}

func NewStore() *Store { return &Store{} }

func (s *Store) Begin() *Transaction {
	return &Transaction{store: s}
}

func (s *Store) Snapshot() []domain.Batch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Batch, len(s.batches))
	for i, batch := range s.batches {
		out[i] = batch.Clone()
	}
	return out
}

type Transaction struct {
	mu     sync.Mutex
	store  *Store
	staged []domain.Batch
	closed bool
}

func (t *Transaction) Stage(ctx context.Context, batch domain.Batch) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := batch.Validate(); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return ErrTransactionClosed
	}
	t.staged = append(t.staged, batch.Clone())
	return nil
}

func (t *Transaction) Commit(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return ErrTransactionClosed
	}
	t.store.mu.Lock()
	for _, batch := range t.staged {
		t.store.batches = append(t.store.batches, batch.Clone())
	}
	t.store.mu.Unlock()
	t.staged = nil
	t.closed = true
	return nil
}

func (t *Transaction) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.staged = nil
	t.closed = true
	return nil
}

func (t *Transaction) StagedCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.staged)
}
