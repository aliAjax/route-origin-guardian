package adapter

import (
	"context"
	"errors"
	"net/http"

	validation "github.com/routeorigin/route-origin-guardian/internal/validation/domain"
)

func HTTPStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case failureMatches(err, context.Canceled):
		return 499
	case failureMatches(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout
	case failureMatches(err, validation.ErrNoCoveringROA):
		return http.StatusNotFound
	case failureMatches(err, validation.ErrOriginMismatch):
		return http.StatusUnprocessableEntity
	case failureMatches(err, validation.ErrSnapshotUnavailable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func failureMatches(err, target error) bool { return errors.Is(err, target) }
