package lifecycle

import "context"

func WaitForExit(ctx context.Context, _ <-chan error) error {
	<-ctx.Done()
	return ctx.Err()
}
