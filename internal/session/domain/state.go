package domain

import "time"

type State string

const (
	Idle        State = "idle"
	OpenSent    State = "open-sent"
	Established State = "established"
	Closing     State = "closing"
	Failed      State = "failed"
)

type Session struct {
	ID          string
	RouterID    string
	PeerAddress string
	RemoteASN   uint32
	State       State
	LastMessage time.Time
	Hold        time.Duration
	RetryAt     time.Time
	Attempts    int
}

func (s Session) Healthy(now time.Time) bool {
	return s.State == Established && s.Hold > 0 && now.Before(s.LastMessage.Add(s.Hold))
}
func (s *Session) Open(now time.Time) error {
	switch s.State {
	case Idle:
	case Failed:
		if now.Before(s.RetryAt) {
			return ErrRetryNotReady
		}
	default:
		return ErrTransition
	}
	s.State = OpenSent
	s.LastMessage = now
	return nil
}
func (s *Session) Established(now time.Time) error {
	if s.State != OpenSent {
		return ErrTransition
	}
	s.State = Established
	s.LastMessage = now
	s.Attempts = 0
	s.RetryAt = time.Time{}
	return nil
}
func (s *Session) Fail(now time.Time) {
	s.State = Failed
	s.LastMessage = now
	s.Attempts++
	s.RetryAt = now.Add(backoff(s.Attempts))
}
func (s Session) ReadyToRetry(now time.Time) bool {
	return s.State == Failed && !now.Before(s.RetryAt)
}
func backoff(n int) time.Duration {
	if n < 1 {
		n = 1
	}
	if n > 6 {
		n = 6
	}
	return time.Duration(1<<(n-1)) * time.Second
}

var (
	ErrTransition    = errorString("invalid session transition")
	ErrRetryNotReady = errorString("session retry is not ready")
)

type errorString string

func (e errorString) Error() string { return string(e) }
