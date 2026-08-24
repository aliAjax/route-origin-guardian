package application

import (
	"context"
	"fmt"
)

func WaitAdmission(ctx context.Context, ready <-chan struct{}) error {
	if ctx == nil || ready == nil {
		return fmt.Errorf("BMP admission wait inputs are required")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ready:
		return nil
	}
}
