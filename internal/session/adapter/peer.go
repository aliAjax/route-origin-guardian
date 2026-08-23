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
	return now.Add(k.Interval)
}
