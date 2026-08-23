package storage

import (
	"context"
	"sync"
	"time"
)

type Checkpoint struct {
	Observer  string
	Sequence  uint64
	UpdatedAt time.Time
}
type Checkpoints struct {
	mu    sync.RWMutex
	items map[string]Checkpoint
}

func NewCheckpoints() *Checkpoints { return &Checkpoints{items: make(map[string]Checkpoint)} }
func (c *Checkpoints) Save(_ context.Context, v Checkpoint) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if old, ok := c.items[v.Observer]; ok && v.Sequence < old.Sequence {
		return nil
	}
	c.items[v.Observer] = v
	return nil
}
func (c *Checkpoints) Get(_ context.Context, id string) (Checkpoint, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.items[id]
	return v, ok
}
