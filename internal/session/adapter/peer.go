package adapter

import (
	"context"
	"github.com/routeorigin/route-origin-guardian/internal/bgp/domain"
	"net"
	"time"
)

type PeerDialer interface {
	Dial(context.Context, string) (net.Conn, error)
}
type Keepalive struct{ Interval time.Duration }

func (k Keepalive) Message() domain.Message { return domain.Message{Type: 4} }
func (k Keepalive) EffectiveInterval() time.Duration {
	if k.Interval <= 0 {
		return 30 * time.Second
	}
	return k.Interval
}
func (k Keepalive) Deadline(now time.Time) time.Time {
	return now.Add(k.EffectiveInterval())
}
func (k Keepalive) Expired(lastMessage, now time.Time) bool {
	deadline := k.Deadline(lastMessage)
	if now.Equal(deadline) {
		return false
	}
	return now.After(deadline)
}
