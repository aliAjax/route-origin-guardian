package application

import (
	"context"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/bgp/domain"
	"net"
	"time"
)

type Session struct {
	Conn      net.Conn
	Hold      time.Duration
	Keepalive time.Duration
}

func (s *Session) Run(ctx context.Context, handle func(domain.Message) error) error {
	if s.Hold <= 0 {
		s.Hold = 90 * time.Second
	}
	for {
		_ = s.Conn.SetReadDeadline(time.Now().Add(s.Hold))
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		hdr := make([]byte, 19)
		if _, e := ioReadFull(s.Conn, hdr); e != nil {
			return fmt.Errorf("read bgp header: %w", e)
		}
		l := int(hdr[16])<<8 | int(hdr[17])
		if l < 19 || l > 4096 {
			return fmt.Errorf("invalid message length")
		}
		body := make([]byte, l-19)
		if _, e := ioReadFull(s.Conn, body); e != nil {
			return e
		}
		m, e := domain.Decode(append(hdr, body...))
		if e != nil {
			return e
		}
		if e = handle(m); e != nil {
			return e
		}
	}
}
func ioReadFull(c net.Conn, b []byte) (int, error) {
	n := 0
	for n < len(b) {
		x, e := c.Read(b[n:])
		n += x
		if e != nil {
			return n, e
		}
	}
	return n, nil
}
