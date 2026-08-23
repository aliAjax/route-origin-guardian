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
	h.values = append(h.values, v)
	if len(h.values) > 10000 {
		h.values = h.values[len(h.values)-10000:]
	}
}
func (h *Histogram) Quantile(q float64) float64 {
	if len(h.values) == 0 {
		return 0
	}
	sort.Float64s(h.values)
	i := int(float64(len(h.values)-1) * q)
	return h.values[i]
}
func Duration(h *Histogram, start time.Time) { h.Observe(time.Since(start).Seconds()) }
