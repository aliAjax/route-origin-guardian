package domain_test

import (
	"errors"
	"testing"
	"time"

	app "github.com/routeorigin/route-origin-guardian/internal/analysis/application"
	"github.com/routeorigin/route-origin-guardian/internal/analysis/domain"
	"github.com/routeorigin/route-origin-guardian/internal/analysis/infrastructure"
	rib "github.com/routeorigin/route-origin-guardian/internal/rib/domain"
)

func TestPolicyPreservesRetryableError(t *testing.T) {
	route := rib.Route{ID: "r-1"}
	if !domain.IsRetryable(domain.ErrRetryable) {
		t.Fatal("retryable classifier rejected sentinel")
	}
	if got := domain.RetryableKind(domain.ErrRetryable); got != "retryable" {
		t.Fatalf("retryable kind = %q", got)
	}
	if err := app.EvaluateWithError(app.Policy{}, route); !errors.Is(err, domain.ErrRetryable) {
		t.Fatalf("policy lost retryable sentinel: %v", err)
	}
	s := infrastructure.NewSuppressor(time.Minute)
	if err := s.RecordError("r-1", time.Unix(10, 0), domain.ErrRetryable); !errors.Is(err, domain.ErrRetryable) {
		t.Fatalf("suppressor lost retryable sentinel: %v", err)
	}
	if err := s.RecordError("r-1", time.Unix(11, 0), domain.ErrRetryable); !errors.Is(err, domain.ErrRetryable) {
		t.Fatalf("suppressed retryable sentinel lost: %v", err)
	}
}
