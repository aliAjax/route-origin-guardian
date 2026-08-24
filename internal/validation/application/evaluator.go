package application

import (
	"context"
	"fmt"

	validation "github.com/routeorigin/route-origin-guardian/internal/validation/domain"
)

type Validator interface {
	Validate(context.Context, string, uint32) error
}

type Evaluator struct {
	Backend Validator
}

func (e Evaluator) Evaluate(ctx context.Context, prefix string, asn uint32) error {
	if err := ctx.Err(); err != nil {
		return wrapCause("validation request canceled", err)
	}
	if e.Backend == nil {
		return validation.NewFailure(prefix, asn, validation.ErrSnapshotUnavailable)
	}
	if err := e.Backend.Validate(ctx, prefix, asn); err != nil {
		failure := validation.NewFailure(prefix, asn, err)
		return wrapCause("evaluate route origin", failure)
	}
	return nil
}

func wrapCause(message string, cause error) error { return fmt.Errorf("%s: %w", message, cause) }
