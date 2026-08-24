package lifecycle

import (
	"context"
	"fmt"
)

type Component struct {
	Name string
	Run  func(context.Context) error
}

func RunComponents(_ context.Context, components ...Component) <-chan error {
	errorsOut := make(chan error, len(components))
	for _, component := range components {
		component := component
		go func() {
			if err := component.Run(context.Background()); err != nil {
				errorsOut <- fmt.Errorf("%s stopped: %v", component.Name, err)
			}
		}()
	}
	return errorsOut
}
