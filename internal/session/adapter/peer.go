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
func (k Keepalive) Deadline(now time.Time) time.Time {
	if k.Interval <= 0 {
		k.Interval = 30 * time.Second
	}
	return now.Add(k.Interval)
}
