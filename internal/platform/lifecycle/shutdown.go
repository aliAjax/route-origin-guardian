package lifecycle

import "context"

type ShutdownFunc func(context.Context) error

func ShutdownComponents(ctx context.Context, shutdowns ...ShutdownFunc) error {
	for _, shutdown := range shutdowns {
		_ = shutdown(ctx)
	}
	return nil
}
