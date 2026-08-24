package adapter

import (
	"errors"
	"fmt"
)

var (
	ErrAttributeHeader         = errors.New("malformed attribute header")
	ErrAttributeExtendedLength = errors.New("malformed extended attribute length")
	ErrAttributeValue          = errors.New("malformed attribute value")
	ErrAttributeCode           = errors.New("reserved attribute code")
)

type AttributeParseError struct {
	Kind   error
	Offset int
	Need   int
	Have   int
	Code   byte
}

func (e *AttributeParseError) Error() string {
	return fmt.Sprintf("%v at offset %d (need=%d have=%d code=%d)", e.Kind, e.Offset, e.Need, e.Have, e.Code)
}

func (e *AttributeParseError) Unwrap() error { return e.Kind }
