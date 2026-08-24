package application_test

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

type rejectingValidator struct {
	err error
}

func (v rejectingValidator) Validate(context.Context, string, uint32) error {
	return v.err
}

func TestValidationFailurePreservesDomainCause(t *testing.T) {
	err := validation.NewFailure("203.0.113.0/24", 64500, validation.ErrNoCoveringROA)
	if !errors.Is(err, validation.ErrNoCoveringROA) {
		t.Fatalf("domain failure lost no-covering-ROA cause: %v", err)
	}
}

func TestValidationEvaluatorPreservesBackendCause(t *testing.T) {
	evaluator := validationapp.Evaluator{Backend: rejectingValidator{err: validation.ErrOriginMismatch}}
	err := evaluator.Evaluate(context.Background(), "203.0.113.0/24", 64599)
	if !errors.Is(err, validation.ErrOriginMismatch) {
		t.Fatalf("evaluator lost origin-mismatch cause: %v", err)
	}
}

func TestValidationFailureStoreRoundTrip(t *testing.T) {
	evaluator := validationapp.Evaluator{Backend: rejectingValidator{err: validation.ErrSnapshotUnavailable}}
	err := evaluator.Evaluate(context.Background(), "198.51.100.0/24", 64501)
	record, encodeErr := validationstore.EncodeFailure(err)
	if encodeErr != nil {
		t.Fatalf("encode failure: %v", encodeErr)
	}
	restored := validationstore.DecodeFailure(record)
	if !errors.Is(restored, validation.ErrSnapshotUnavailable) {
		t.Fatalf("restored failure lost snapshot-unavailable cause: %v", restored)
	}
}

func TestValidationHTTPStatusUsesWrappedCause(t *testing.T) {
	evaluator := validationapp.Evaluator{Backend: rejectingValidator{err: validation.ErrOriginMismatch}}
	err := evaluator.Evaluate(context.Background(), "192.0.2.0/24", 64511)
	if got := validationadapter.HTTPStatus(err); got != http.StatusUnprocessableEntity {
		t.Fatalf("HTTPStatus()=%d, want %d for wrapped origin mismatch", got, http.StatusUnprocessableEntity)
	}
}
