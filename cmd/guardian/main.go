package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/routeorigin/route-origin-guardian/internal/bmp/adapter"
	"github.com/routeorigin/route-origin-guardian/internal/platform/config"
	platformhttp "github.com/routeorigin/route-origin-guardian/internal/platform/httpx"
	"github.com/routeorigin/route-origin-guardian/internal/platform/lifecycle"
	"github.com/routeorigin/route-origin-guardian/internal/platform/logging"
	"github.com/routeorigin/route-origin-guardian/internal/rib/application"
	"github.com/routeorigin/route-origin-guardian/internal/rib/infrastructure"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/domain"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}
	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)
	store := infrastructure.NewMemoryStore()
	service := application.NewService(store)
	roas := domain.NewIndex()
	server := platformhttp.NewServer(cfg, service, roas)
	collector := adapter.NewCollector(service, cfg.BMPListen, cfg.BGPListen)
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx := lifecycle.RuntimeContext(signalCtx)
	runtimeErrors := lifecycle.RunComponents(ctx,
		lifecycle.Component{Name: "collector", Run: collector.Run},
		lifecycle.Component{Name: "http", Run: func(context.Context) error { return server.ListenAndServe() }},
	)
	if err := lifecycle.WaitForExit(ctx, runtimeErrors); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("runtime stopped", "error", err)
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := lifecycle.ShutdownComponents(shutdown, server.Shutdown, func(context.Context) error { return collector.Close() }); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
	logger.Info("route origin guardian stopped")
}

// keep net imported in the binary for operational probes and future listener wiring.
var _ = net.IPv4len
