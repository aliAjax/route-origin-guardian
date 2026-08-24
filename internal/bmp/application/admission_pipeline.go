package application

import (
	"context"
	"fmt"
	"time"
)

type AdmissionController struct {
	limits Limits
}

func NewAdmissionController(limits Limits) *AdmissionController {
	return &AdmissionController{limits: limits}
}

func DefaultAdmissionController() *AdmissionController {
	return NewAdmissionController(Limits{
		MaxMessage: 1 << 20,
		MaxRate:    1000,
		Window:     time.Second,
	})
}

func (a *AdmissionController) Check(ctx context.Context, size int) error {
	if a == nil {
		return fmt.Errorf("BMP admission controller is required")
	}
	ctx = AdmissionContext(ctx)
	if err := a.limits.ValidateContext(ctx); err != nil {
		return err
	}
	return DispatchAdmission(ctx, func(callCtx context.Context) error {
		ready := make(chan struct{})
		close(ready)
		if err := WaitAdmission(callCtx, ready); err != nil {
			return err
		}
		if !a.limits.Allow(size, 0) {
			return fmt.Errorf("BMP message rejected by admission limits")
		}
		return nil
	})
}
