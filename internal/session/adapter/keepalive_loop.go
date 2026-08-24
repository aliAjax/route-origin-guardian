package adapter

import (
	"context"
	"net"
	"time"
)

func RunKeepalives(_ context.Context, conn net.Conn, interval time.Duration) error {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if _, err := conn.Write([]byte{4}); err != nil {
				return err
			}
		}
	}
}
