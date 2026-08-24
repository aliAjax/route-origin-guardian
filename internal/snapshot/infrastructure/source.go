package infrastructure

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/routeorigin/route-origin-guardian/internal/snapshot/domain"
)

var ErrSubscriptionClosed = errors.New("snapshot subscription closed")

type Source struct {
	batches  []domain.Batch
	closeErr error
}

func NewSource(batches []domain.Batch, closeErr error) *Source {
	cloned := make([]domain.Batch, len(batches))
	for i, batch := range batches {
		cloned[i] = batch.Clone()
	}
	return &Source{batches: cloned, closeErr: closeErr}
}

func (s *Source) Open(ctx context.Context) (*Subscription, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &Subscription{batches: s.batches, closeErr: s.closeErr}, nil
}

type Subscription struct {
	mu         sync.Mutex
	batches    []domain.Batch
	index      int
	closed     bool
	closeCount int
	closeErr   error
}

func (s *Subscription) Next(ctx context.Context) (domain.Batch, error) {
	if err := ctx.Err(); err != nil {
		return domain.Batch{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.Batch{}, ErrSubscriptionClosed
	}
	if s.index >= len(s.batches) {
		return domain.Batch{}, io.EOF
	}
	batch := s.batches[s.index].Clone()
	s.index++
	return batch, nil
}

func (s *Subscription) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	s.closeCount++
	return s.closeErr
}

func (s *Subscription) CloseCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closeCount
}
