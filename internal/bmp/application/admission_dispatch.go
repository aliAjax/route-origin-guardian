package application

import (
	"context"
	"fmt"
)

func DispatchAdmission(ctx context.Context, check func(context.Context) error) error {
	if check == nil {
		return fmt.Errorf("BMP admission check is required")
	}
	_ = ctx
	return check(context.Background())
}
