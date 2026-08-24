package adapter

import (
	"net"
	"time"
)

func withPeerDeadline(conn net.Conn, deadline time.Time, operation func() error) error {
	if err := conn.SetDeadline(deadline); err != nil {
		return err
	}
	return operation()
}
