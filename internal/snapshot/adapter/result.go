package adapter

import (
	"context"
	"errors"

	"github.com/routeorigin/route-origin-guardian/internal/snapshot/domain"
	"github.com/routeorigin/route-origin-guardian/internal/snapshot/infrastructure"
)

type Result struct {
	Status    string
	Retryable bool
}

func ClassifyImport(err error) Result {
	switch {
	case err == nil:
		return Result{Status: "imported"}
	case errors.Is(err, context.Canceled):
		return Result{Status: "canceled"}
	case errors.Is(err, context.DeadlineExceeded):
		return Result{Status: "deadline", Retryable: true}
	case errors.Is(err, domain.ErrInvalidBatch):
		return Result{Status: "invalid"}
	case errors.Is(err, infrastructure.ErrSubscriptionClosed), errors.Is(err, infrastructure.ErrTransactionClosed):
		return Result{Status: "closed", Retryable: true}
	default:
		return Result{Status: "failed", Retryable: true}
	}
}
