package application

import "time"

func (q *Quota) advanceWindow(now time.Time) {
	if now.Sub(q.start) >= q.limits.Window {
		q.start = now
		q.used = 0
	}
}
