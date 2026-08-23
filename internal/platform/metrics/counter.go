package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct{ value atomic.Uint64 }

func (c *Counter) Inc()         { c.value.Add(1) }
func (c *Counter) Add(n uint64) { c.value.Add(n) }
func (c *Counter) Get() uint64  { return c.value.Load() }

type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

func NewRegistry() *Registry { return &Registry{counters: make(map[string]*Counter)} }
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	// A zero-value Registry (e.g. an embedded struct field that was never
	// constructed via NewRegistry) has a nil counters map; lazily allocate it
	// so the first counter registration does not panic with
	// "assignment to entry in nil map".
	if r.counters == nil {
		r.counters = make(map[string]*Counter)
	}
	if c := r.counters[name]; c != nil {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}
func (r *Registry) Text() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := ""
	for n, c := range r.counters {
		out += fmt.Sprintf("%s %d\n", n, c.Get())
	}
	return out
}
