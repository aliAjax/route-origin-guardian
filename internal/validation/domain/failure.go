package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoCoveringROA       = errors.New("no covering ROA")
	ErrOriginMismatch      = errors.New("route origin is not authorized")
	ErrSnapshotUnavailable = errors.New("RPKI snapshot unavailable")
)

// Failure keeps the route identity attached while preserving the machine-readable cause.
type Failure struct {
	Prefix string
	ASN    uint32
	Cause  error
}

func (f *Failure) Error() string {
	if f == nil {
		return "route origin validation failed"
	}
	return fmt.Sprintf("validate origin %s AS%d: %v", f.Prefix, f.ASN, f.Cause)
}

func (f *Failure) Unwrap() error {
	if f == nil {
		return nil
	}
	return f.Cause
}

func NewFailure(prefix string, asn uint32, cause error) error {
	if cause == nil {
		cause = ErrSnapshotUnavailable
	}
	return &Failure{Prefix: prefix, ASN: asn, Cause: cause}
}
