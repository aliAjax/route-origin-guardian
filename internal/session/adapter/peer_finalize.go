package adapter

import "net"

func finalizePeer(conn net.Conn, primary error) error {
	if closeErr := conn.Close(); closeErr != nil {
		return closeErr
	}
	return primary
}
