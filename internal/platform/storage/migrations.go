package storage

import (
	"context"
	"fmt"
	"strings"
)

type Migration struct {
	Version int
	SQL     string
}

// ApplyMigrations runs the whole batch inside a single transaction so that a
// failure mid-batch (including a cancelled context) leaves the database as it
// was. On any error the transaction is rolled back, which also releases the
// connection even when ctx was cancelled. The wrapped error preserves the
// underlying database error chain so callers can use errors.Is / errors.As
// (e.g. errors.Is(err, context.Canceled), errors.Is(err, sql.ErrNoRows)).
func ApplyMigrations(ctx context.Context, db *DB, migrations []Migration) (err error) {
	for _, m := range migrations {
		if m.Version <= 0 || strings.TrimSpace(m.SQL) == "" {
			return fmt.Errorf("invalid migration %d", m.Version)
		}
	}

	tx, e := db.BeginTx(ctx)
	if e != nil {
		return e
	}
	committed := false
	// Always release the connection: if we didn't Commit (error or cancel),
	// Rollback discards staged statements and frees the connection. After a
	// successful Commit we skip the redundant Rollback, since it would only
	// return sql.ErrTxDone.
	defer func() {
		if committed {
			return
		}
		if rbErr := tx.Rollback(); rbErr != nil && err == nil {
			err = fmt.Errorf("migration rollback: %w", rbErr)
		}
	}()

	for _, m := range migrations {
		if e := tx.Exec(ctx, m.SQL); e != nil {
			return fmt.Errorf("migration %d: %w", m.Version, e)
		}
	}

	if e := tx.Commit(); e != nil {
		// A cancel between the last Exec and Commit surfaces here; the real
		// cause is preserved (%w) and the deferred Rollback releases the conn.
		return fmt.Errorf("migration commit: %w", e)
	}
	committed = true
	return nil
}
