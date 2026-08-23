package metrics

import (
	"sort"
	"sync"
	"time"
)

type Histogram struct {
	mu     sync.Mutex
	values []float64
}

func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	h.values = append(h.values, v)
	if len(h.values) > 10000 {
		h.values = h.values[len(h.values)-10000:]
	}
	h.mu.Unlock()
}
func (h *Histogram) Quantile(q float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.values) == 0 {
		return 0
	}
	cp := append([]float64(nil), h.values...)
	sort.Float64s(cp)
	i := int(float64(len(cp)-1) * q)
	return cp[i]
}
func Duration(h *Histogram, start time.Time) { h.Observe(time.Since(start).Seconds()) }
