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

func (s JSONSource) Fetch(_ context.Context) (application.Snapshot, error) {
	var snap application.Snapshot
	if s.Reader == nil {
		return snap, fmt.Errorf("nil ROA reader")
	}
	if e := json.NewDecoder(s.Reader).Decode(&snap); e != nil {
		return snap, fmt.Errorf("decode ROA snapshot: %w", e)
	}
	return snap, nil
}
