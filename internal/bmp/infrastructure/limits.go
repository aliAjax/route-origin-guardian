package infrastructure

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu           sync.Mutex
	window       time.Time
	count, limit int
}

func NewRateLimiter(limit int) *RateLimiter { return &RateLimiter{limit: limit} }
func (l *RateLimiter) Allow(now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.window.IsZero() || now.Sub(l.window) >= time.Second {
		l.window = now
		l.count = 0
	}
	if l.count >= l.limit {
		return false
	}
	l.count++
	return true
}

type ByteBudget struct {
	mu        sync.Mutex
	used, max int64
}

func NewByteBudget(max int64) *ByteBudget { return &ByteBudget{max: max} }
func (b *ByteBudget) Reserve(n int64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if n < 0 || b.used+n > b.max {
		return false
	}
	b.used += n
	return true
}
func (b *ByteBudget) Release(n int64) {
	b.mu.Lock()
	b.used -= n
	if b.used < 0 {
		b.used = 0
	}
	b.mu.Unlock()
}
