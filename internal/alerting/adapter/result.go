package adapter

import (
	"errors"
)

var ErrPermanentDelivery = errors.New("permanent alert delivery failure")

type DeliveryError struct {
	Cause error
}

func (e *DeliveryError) Error() string {
	if e == nil || e.Cause == nil {
		return "alert delivery failed"
	}
	return "alert delivery failed: " + e.Cause.Error()
}

func (e *DeliveryError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

type Result struct {
	Status    string
	Retryable bool
}

func Classify(err error) Result {
	if err == nil {
		return Result{Status: "delivered"}
	}
	if errors.Is(err, ErrPermanentDelivery) {
		return Result{Status: "rejected"}
	}
	return Result{Status: "failed", Retryable: true}
}
