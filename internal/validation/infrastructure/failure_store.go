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
	failure, ok := err.(*validation.Failure)
	if !ok {
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
	var detail string
	switch record.Code {
	case "no_covering_roa":
		detail = validation.ErrNoCoveringROA.Error()
	case "origin_mismatch":
		detail = validation.ErrOriginMismatch.Error()
	case "snapshot_unavailable":
		detail = validation.ErrSnapshotUnavailable.Error()
	default:
		return fmt.Errorf("decode validation failure: unknown code %q", record.Code)
	}
	return validation.NewFailure(record.Prefix, record.ASN, errors.New(detail))
}
