package openapi

import (
	"context"
	"sync"
	"testing"

	ribapp "github.com/routeorigin/route-origin-guardian/internal/rib/application"
	ribdomain "github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	ribinfra "github.com/routeorigin/route-origin-guardian/internal/rib/infrastructure"
)

func TestRouteSnapshotIsolationAndAbort(t *testing.T) {
	store := ribinfra.NewMemoryStore()
	route := ribdomain.Route{TenantID: "tenant", ObserverID: "obs", PeerAddress: "192.0.2.1", Prefix: "10.0.0.0/24", ASPath: []uint32{64512}}
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); <-start; _, _ = store.Upsert(context.Background(), route) }()
	go func() {
		defer wg.Done()
		<-start
		_, _ = store.List(context.Background(), ribapp.Query{Limit: 100})
	}()
	close(start)
	wg.Wait()
	route.ASPath[0] = 64513
	rows, err := store.List(context.Background(), ribapp.Query{Limit: 100})
	if err != nil || len(rows) != 1 || rows[0].ASPath[0] != 64512 {
		t.Fatalf("stored route was aliased: rows=%v err=%v", rows, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ribapp.NewService(store).Announce(ctx, route); err != context.Canceled {
		t.Fatalf("cancelled upsert error=%v", err)
	}
	if _, err := ribapp.NewService(store).Replace(ctx, route); err != context.Canceled {
		t.Fatalf("cancelled replace error=%v", err)
	}
}
