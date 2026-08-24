package infrastructure

import (
	"errors"
	"fmt"

	validation "github.com/routeorigin/route-origin-guardian/internal/validation/domain"
)

type FailureRecord struct {
	Prefix string
	ASN    uint32
	Code   string
	Detail string
}

func EncodeFailure(err error) (FailureRecord, error) {
	var failure *validation.Failure
	if !errors.As(err, &failure) {
		return FailureRecord{}, fmt.Errorf("encode validation failure: %w", err)
	}

	code := ""
	switch {
	case errors.Is(err, validation.ErrNoCoveringROA):
		code = "no_covering_roa"
	case errors.Is(err, validation.ErrOriginMismatch):
		code = "origin_mismatch"
	case errors.Is(err, validation.ErrSnapshotUnavailable):
		code = "snapshot_unavailable"
	default:
		return FailureRecord{}, fmt.Errorf("encode validation failure: unsupported cause: %w", err)
	}
	return FailureRecord{Prefix: failure.Prefix, ASN: failure.ASN, Code: code, Detail: err.Error()}, nil
}

func DecodeFailure(record FailureRecord) error {
	var cause error
	switch record.Code {
	case "no_covering_roa":
		cause = validation.ErrNoCoveringROA
	case "origin_mismatch":
		cause = validation.ErrOriginMismatch
	case "snapshot_unavailable":
		cause = validation.ErrSnapshotUnavailable
	default:
		return fmt.Errorf("decode validation failure: unknown code %q", record.Code)
	}
	return validation.NewFailure(record.Prefix, record.ASN, cause)
}
