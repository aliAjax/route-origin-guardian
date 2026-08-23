package adapter

import (
	"context"
	"net"
)

func bindContext(_ context.Context, conn net.Conn) net.Conn {
	return conn
}
