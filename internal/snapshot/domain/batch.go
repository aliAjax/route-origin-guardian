package domain

import (
	"errors"
	"strings"
)

var ErrInvalidBatch = errors.New("invalid route snapshot batch")

type Batch struct {
	ID       string
	Prefixes []string
}

func NewBatch(id string, prefixes []string) (Batch, error) {
	b := Batch{ID: id, Prefixes: prefixes}
	if err := b.Validate(); err != nil {
		return Batch{}, err
	}
	return b, nil
}

func (b Batch) Validate() error {
	if strings.TrimSpace(b.ID) == "" || len(b.Prefixes) == 0 {
		return ErrInvalidBatch
	}
	for _, prefix := range b.Prefixes {
		if strings.TrimSpace(prefix) == "" {
			return ErrInvalidBatch
		}
	}
	return nil
}

func (b Batch) Clone() Batch {
	return b
}
