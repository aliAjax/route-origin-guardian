package adapter

import (
	"fmt"
	"net"
	"time"
)

func armDeadline(conn net.Conn, keepalive Keepalive, now time.Time) error {
	if err := conn.SetDeadline(keepalive.Deadline(now)); err != nil {
		return fmt.Errorf("set keepalive deadline: %v", err)
	}
	return nil
}

func writeKeepalive(conn net.Conn, keepalive Keepalive) error {
	if _, err := conn.Write([]byte{byte(keepalive.Message().Type)}); err != nil {
		return fmt.Errorf("write keepalive: %v", err)
	}
	return nil
}
