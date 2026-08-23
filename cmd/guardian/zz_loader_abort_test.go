package main

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/routeorigin/route-origin-guardian/internal/rpki/adapter"
	app "github.com/routeorigin/route-origin-guardian/internal/rpki/application"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/domain"
	infra "github.com/routeorigin/route-origin-guardian/internal/rpki/infrastructure"
)

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("reader reached") }

type failingSource struct{}

func (failingSource) Fetch(context.Context) (app.Snapshot, error) {
	return app.Snapshot{}, errors.New("source reached")
}

type cancelingReader struct {
	data   *strings.Reader
	cancel context.CancelFunc
}

func (r *cancelingReader) Read(p []byte) (int, error) {
	n, err := r.data.Read(p)
	if n > 0 {
		r.cancel()
	}
	return n, err
}

var _ io.Reader = (*cancelingReader)(nil)

func TestLoaderStopsOnContextAbort(t *testing.T) {
	index := domain.NewIndex()
	service := app.NewService(index)
	loader := infra.NewLoader(failingSource{}, service, time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := loader.Load(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	if index.Version() != 0 {
		t.Fatalf("cancelled load changed index version: %d", index.Version())
	}
	preCtx, preCancel := context.WithCancel(context.Background())
	preCancel()
	if _, err := (adapter.JSONSource{Reader: failingReader{}}).Fetch(preCtx); !errors.Is(err, context.Canceled) {
		t.Fatalf("source ignored pre-cancel: %v", err)
	}
	readCtx, readCancel := context.WithCancel(context.Background())
	reader := adapter.JSONSource{Reader: &cancelingReader{data: strings.NewReader(`{"serial":9,"roas":[]}`), cancel: readCancel}}
	if _, err := reader.Fetch(readCtx); !errors.Is(err, context.Canceled) {
		t.Fatalf("source ignored cancellation during read: %v", err)
	}
	if err := service.Import(ctx, app.Snapshot{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("direct import ignored cancellation: %v", err)
	}
}
