package infrastructure

import (
	"fmt"
	"sync"
	"time"
)

type Suppressor struct {
	mu     sync.Mutex
	until  map[string]time.Time
	window time.Duration
}

func NewSuppressor(window time.Duration) *Suppressor {
	if window <= 0 {
		window = time.Minute
	}
	return &Suppressor{until: make(map[string]time.Time), window: window}
}
func (s *Suppressor) Allow(key string, now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t, ok := s.until[key]; ok && now.Before(t) {
		return false
	}
	s.until[key] = now.Add(s.window)
	return true
}
func (s *Suppressor) Clear(key string) { s.mu.Lock(); delete(s.until, key); s.mu.Unlock() }

func (s *Suppressor) RecordError(key string, now time.Time, err error) error {
	if !s.Allow(key, now) {
		return fmt.Errorf("suppressed: %w", err)
	}
	return fmt.Errorf("analysis failed: %w", err)
}
