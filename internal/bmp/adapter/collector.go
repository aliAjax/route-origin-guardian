package adapter

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rib/application"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"net"
	"strings"
	"sync"
)

// Collector accepts a bounded JSON line representation for local BMP/BGP smoke tests.
// Production deployments can replace it with protocol adapters implementing the same service calls.
type Collector struct {
	svc              *application.Service
	bmpAddr, bgpAddr string
	listeners        []net.Listener
	mu               sync.Mutex
}

func NewCollector(s *application.Service, bmp, bgp string) *Collector {
	return &Collector{svc: s, bmpAddr: bmp, bgpAddr: bgp}
}
func (c *Collector) Run(ctx context.Context) error {
	for _, a := range []string{c.bmpAddr, c.bgpAddr} {
		ln, e := net.Listen("tcp", a)
		if e != nil {
			return fmt.Errorf("listen %s: %w", a, e)
		}
		c.mu.Lock()
		c.listeners = append(c.listeners, ln)
		c.mu.Unlock()
		go c.accept(ctx, ln)
	}
	<-ctx.Done()
	return ctx.Err()
}
func (c *Collector) accept(ctx context.Context, ln net.Listener) {
	for {
		conn, e := ln.Accept()
		if e != nil {
			return
		}
		go c.handle(ctx, conn)
	}
}
func (c *Collector) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	scan := bufio.NewScanner(ioLimitReader{r: conn, n: 1 << 20})
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		var in struct {
			Operation string       `json:"operation"`
			Route     domain.Route `json:"route"`
		}
		if e := json.Unmarshal([]byte(line), &in); e != nil {
			fmt.Fprintf(conn, "{\"error\":%q}\n", e.Error())
			continue
		}
		var e error
		switch strings.ToLower(in.Operation) {
		case "withdraw":
			_, e = c.svc.Withdraw(ctx, in.Route)
		default:
			_, e = c.svc.Announce(ctx, in.Route)
		}
		if e != nil {
			fmt.Fprintf(conn, "{\"error\":%q}\n", e.Error())
			continue
		}
		fmt.Fprintln(conn, "{\"status\":\"accepted\"}")
	}
}
func (c *Collector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, l := range c.listeners {
		_ = l.Close()
	}
	c.listeners = nil
	return nil
}

type ioLimitReader struct {
	r net.Conn
	n int64
}

func (i ioLimitReader) Read(p []byte) (int, error) {
	if i.n <= 0 {
		return 0, fmt.Errorf("message limit exceeded")
	}
	if int64(len(p)) > i.n {
		p = p[:i.n]
	}
	n, e := i.r.Read(p)
	i.n -= int64(n)
	return n, e
}
