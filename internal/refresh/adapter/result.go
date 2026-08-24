package adapter

import (
	"context"
	"errors"
)

type Result struct {
	Status    string
	Retryable bool
}

func Classify(err error) Result {
	switch {
	case err == nil:
		return Result{Status: "complete"}
	case errors.Is(err, context.Canceled):
		return Result{Status: "canceled"}
	case errors.Is(err, context.DeadlineExceeded):
		return Result{Status: "timed_out", Retryable: true}
	default:
		return Result{Status: "failed", Retryable: true}
	}
}
