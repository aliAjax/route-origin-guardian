package infrastructure

import "time"

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

type RetryPolicy struct {
	Base, Max time.Duration
	Attempts  int
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if p.Base <= 0 {
		p.Base = time.Second
	}
	if p.Max <= 0 {
		p.Max = 5 * time.Minute
	}
	if attempt < 1 {
		attempt = 1
	}
	d := p.Base
	for n := 1; n < attempt; n++ {
		if d > p.Max/2 {
			return p.Max
		}
		d *= 2
	}
	if d > p.Max {
		return p.Max
	}
	return d
}
