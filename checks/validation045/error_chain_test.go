package validation045

import (
	"context"
	"errors"
	"net/http"
	"testing"

	validationadapter "github.com/routeorigin/route-origin-guardian/internal/validation/adapter"
	validationapp "github.com/routeorigin/route-origin-guardian/internal/validation/application"
	validation "github.com/routeorigin/route-origin-guardian/internal/validation/domain"
	validationstore "github.com/routeorigin/route-origin-guardian/internal/validation/infrastructure"
)

func TestDomainFailureIdentityR045(t *testing.T) {
	err := validation.NewFailure("203.0.113.0/24", 64500, validation.ErrNoCoveringROA)
	if !errors.Is(err, validation.ErrNoCoveringROA) {
		t.Fatalf("domain failure lost sentinel identity: %v", err)
	}
}

func TestEvaluatorCancellationIdentityR045(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := (validationapp.Evaluator{}).Evaluate(ctx, "203.0.113.0/24", 64500)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("evaluator lost cancellation identity: %v", err)
	}
}

func TestFailureEncodingWrappedCauseR045(t *testing.T) {
	err := validation.NewFailure("198.51.100.0/24", 64501, validation.ErrOriginMismatch)
	err = &outerFailure{cause: err}
	record, encodeErr := validationstore.EncodeFailure(err)
	if encodeErr != nil {
		t.Fatalf("encode wrapped failure: %v", encodeErr)
	}
	if record.Code != "origin_mismatch" {
		t.Fatalf("record code=%q, want origin_mismatch", record.Code)
	}
}

func TestFailureRoundTripIdentityR045(t *testing.T) {
	restored := validationstore.DecodeFailure(validationstore.FailureRecord{
		Prefix: "192.0.2.0/24", ASN: 64511, Code: "snapshot_unavailable",
	})
	if !errors.Is(restored, validation.ErrSnapshotUnavailable) {
		t.Fatalf("decoded failure lost sentinel identity: %v", restored)
	}
}

func TestHTTPStatusWrappedCauseR045(t *testing.T) {
	err := &outerFailure{cause: validation.ErrOriginMismatch}
	if got := validationadapter.HTTPStatus(err); got != http.StatusUnprocessableEntity {
		t.Fatalf("HTTPStatus()=%d, want %d", got, http.StatusUnprocessableEntity)
	}
}

type outerFailure struct {
	cause error
}

func (e *outerFailure) Error() string { return "validation boundary: " + e.cause.Error() }
func (e *outerFailure) Unwrap() error { return e.cause }
