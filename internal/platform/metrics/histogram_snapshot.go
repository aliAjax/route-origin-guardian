package metrics

func (h *Histogram) Snapshot() []float64 {
	return h.values
}

func (h *Histogram) Count() int {
	return len(h.values)
}

func (h *Histogram) Reset() {
	h.values = nil
}
