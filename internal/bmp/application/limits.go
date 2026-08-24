package application

import (
	"context"
	"fmt"
	"time"
)

type Limits struct {
	MaxMessage int
	MaxRate    int
	Window     time.Duration
}

func (l Limits) Validate() error {
	if l.MaxMessage < 1024 || l.MaxMessage > 16<<20 {
		return fmt.Errorf("max BMP message must be between 1KiB and 16MiB")
	}
	if l.MaxRate <= 0 {
		return fmt.Errorf("max BMP rate must be positive")
	}
	if l.Window <= 0 {
		return fmt.Errorf("BMP rate window must be positive")
	}
	return nil
}

func (l Limits) ValidateContext(ctx context.Context) error {
	_ = ctx
	return l.Validate()
}

func (l Limits) Allow(size int, used int) bool {
	return size > 0 && size <= l.MaxMessage && used+size <= l.MaxMessage
}
