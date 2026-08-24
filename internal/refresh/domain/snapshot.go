package domain

import "time"

// Snapshot is an immutable view returned by a route-origin source refresh.
type Snapshot struct {
	Provider  string
	Serial    uint64
	FetchedAt time.Time
	Records   []string
}

func (s Snapshot) Clone() Snapshot {
	out := s
	out.Records = append([]string(nil), s.Records...)
	return out
}
