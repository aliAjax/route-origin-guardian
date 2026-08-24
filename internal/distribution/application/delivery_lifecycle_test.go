package application_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/routeorigin/route-origin-guardian/internal/distribution/application"
	deliverydomain "github.com/routeorigin/route-origin-guardian/internal/distribution/domain"
	"github.com/routeorigin/route-origin-guardian/internal/distribution/infrastructure"
	ribdomain "github.com/routeorigin/route-origin-guardian/internal/rib/domain"
)

type eventSource struct{ events []ribdomain.Event }

func (s eventSource) Events(context.Context, uint64, int) ([]ribdomain.Event, error) {
	return s.events, nil
}

func TestDeliveryLifecycleRetriesThenArchives(t *testing.T) {
	event := ribdomain.Event{ID: 7, Route: ribdomain.Route{ASPath: []uint32{64512}, Communities: []string{"64512:7"}}}
	cursor := application.NewCursor(eventSource{events: []ribdomain.Event{event}})
	items, next, err := cursor.Read(context.Background(), 0, 10)
	if err != nil || next != event.ID || len(items) != 1 {
		t.Fatalf("cursor read: items=%d next=%d err=%v", len(items), next, err)
	}
	event.Route.ASPath[0] = 64513
	if items[0].Route.ASPath[0] != 64512 {
		t.Fatal("cursor result shares route slices with its source")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "retry", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	delivery := deliverydomain.NewDelivery(items[0])
	dead := &infrastructure.DeadLetter{}
	runner := application.DeliveryRunner{
		Sender:      application.WebhookSender{Client: server.Client(), URL: server.URL, Secret: "secret"},
		DeadLetter:  dead,
		MaxAttempts: 2,
	}
	if err := runner.Run(context.Background(), &delivery); err == nil || delivery.Status != deliverydomain.Pending {
		t.Fatalf("first attempt: status=%s err=%v", delivery.Status, err)
	}
	if err := runner.Run(context.Background(), &delivery); err == nil || delivery.Status != deliverydomain.DeadLetter {
		t.Fatalf("second attempt: status=%s err=%v", delivery.Status, err)
	}
	archived := dead.List(context.Background(), 10)
	if len(archived) != 1 || archived[0].ID != event.ID {
		t.Fatalf("dead letters = %+v", archived)
	}
	items[0].Route.Communities[0] = "changed"
	if archived[0].Route.Communities[0] != "64512:7" {
		t.Fatal("dead letter changed with delivery input")
	}
}
