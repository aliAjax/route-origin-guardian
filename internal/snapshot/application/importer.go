package application

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/routeorigin/route-origin-guardian/internal/snapshot/infrastructure"
)

type Importer struct {
	Source *infrastructure.Source
	Store  *infrastructure.Store
}

func (i *Importer) Sync(ctx context.Context) (err error) {
	if i == nil || i.Source == nil || i.Store == nil {
		return fmt.Errorf("snapshot importer unavailable")
	}
	subscription, err := i.Source.Open(ctx)
	if err != nil {
		return fmt.Errorf("open route snapshot: %w", err)
	}
	transaction := i.Store.Begin()
	committed := false
	defer func() {
		if !committed {
			_ = transaction.Rollback()
		}
		err = errors.Join(err, subscription.Close())
	}()

	for {
		batch, nextErr := subscription.Next(ctx)
		if errors.Is(nextErr, io.EOF) {
			break
		}
		if nextErr != nil {
			return fmt.Errorf("read route snapshot: %w", nextErr)
		}
		if stageErr := transaction.Stage(ctx, batch); stageErr != nil {
			return fmt.Errorf("stage route snapshot: %w", stageErr)
		}
	}
	if err := transaction.Commit(ctx); err != nil {
		return fmt.Errorf("commit route snapshot: %w", err)
	}
	committed = true
	return nil
}
