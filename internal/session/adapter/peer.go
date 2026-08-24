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

type HandshakeFunc func(context.Context, net.Conn) error

func DialAndHandshake(ctx context.Context, dialer PeerDialer, address string, timeout time.Duration, handshake HandshakeFunc) (net.Conn, error) {
	conn, err := dialer.Dial(ctx, address)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(timeout)
	if err := withPeerDeadline(conn, deadline, func() error {
		return handshake(ctx, conn)
	}); err != nil {
		return nil, err
	}
	return conn, nil
}

type Keepalive struct{ Interval time.Duration }

func (k Keepalive) Message() domain.Message { return domain.Message{Type: 4} }
func (k Keepalive) Deadline(now time.Time) time.Time {
	if k.Interval <= 0 {
		k.Interval = 30 * time.Second
	}
	return now.Add(k.Interval)
}
