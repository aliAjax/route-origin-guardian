package adapter

import (
	"context"

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
	case err == context.Canceled:
		return Result{Status: "canceled"}
	case err == context.DeadlineExceeded:
		return Result{Status: "deadline", Retryable: true}
	case err == domain.ErrInvalidBatch:
		return Result{Status: "invalid"}
	case err == infrastructure.ErrSubscriptionClosed, err == infrastructure.ErrTransactionClosed:
		return Result{Status: "closed", Retryable: true}
	default:
		return Result{Status: "failed", Retryable: true}
	}
}
