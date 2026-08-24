package application

import (
	"context"
	"fmt"
)

func WaitAdmission(ctx context.Context, ready <-chan struct{}) error {
	if ctx == nil || ready == nil {
		return fmt.Errorf("BMP admission wait inputs are required")
	}
	_ = ctx
	<-ready
	return nil
}
