package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/application"
	"io"
)

type Source interface {
	Fetch(context.Context) (application.Snapshot, error)
}
type JSONSource struct{ Reader io.Reader }

func (s JSONSource) Fetch(ctx context.Context) (application.Snapshot, error) {
	var snap application.Snapshot
	if s.Reader == nil {
		return snap, fmt.Errorf("nil ROA reader")
	}
	if e := ctx.Err(); e != nil {
		return snap, e
	}
	// json.Decoder.Decode cannot be interrupted on an arbitrary io.Reader, so
	// run it on a goroutine and abandon it when the context is cancelled. When
	// the reader is also an io.Closer (e.g. an http.Response.Body) closing it
	// unblocks the in-flight read so the source fetch actually stops instead
	// of running on in the background after the request is gone.
	dec := json.NewDecoder(s.Reader)
	type outcome struct {
		snap application.Snapshot
		err  error
	}
	done := make(chan outcome, 1)
	go func() {
		var o outcome
		o.err = dec.Decode(&o.snap)
		done <- o
	}()
	select {
	case <-ctx.Done():
		if c, ok := s.Reader.(io.Closer); ok {
			_ = c.Close()
		}
		return snap, ctx.Err()
	case o := <-done:
		if o.err != nil {
			return snap, fmt.Errorf("decode ROA snapshot: %w", o.err)
		}
		return o.snap, nil
	}
}
