package adapter

import (
	"context"
	"fmt"
	"net"
	"time"
)

type PeerKeepalive struct {
	Dialer PeerDialer
	Now    func() time.Time
}

func (p PeerKeepalive) Open(ctx context.Context, address string, keepalive Keepalive) (net.Conn, error) {
	conn, err := p.Dialer.Dial(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("dial peer: %v", err)
	}
	if err := armDeadline(conn, keepalive, p.now()); err != nil {
		return nil, err
	}
	conn = bindContext(ctx, conn)
	if err := writeKeepalive(conn, keepalive); err != nil {
		return nil, err
	}
	return conn, nil
}

func (p PeerKeepalive) now() time.Time {
	if p.Now != nil {
		return p.Now()
	}
	return time.Now()
}
