package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type DB struct{ db *sql.DB }

func New(db *sql.DB) *DB { return &DB{db: db} }
func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database unavailable")
	}
	if e := d.db.PingContext(ctx); e != nil {
		return fmt.Errorf("database ping: %w", e)
	}
	return nil
}
func (d *DB) Exec(ctx context.Context, q string, args ...any) error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database unavailable")
	}
	if _, e := d.db.ExecContext(ctx, q, args...); e != nil {
		return fmt.Errorf("database exec: %w", e)
	}
	return nil
}

// Tx is a thin wrapper around sql.Tx that exposes the same Exec signature as
// DB so callers (such as ApplyMigrations) can run a batch of statements inside
// a single transaction and have the underlying connection released on cancel.
type Tx struct{ tx *sql.Tx }

// BeginTx starts a transaction. The ctx governs the Begin call itself; each
// statement is cancelled through the ctx passed to (*Tx).Exec. Wrapping the
// whole batch in one tx is what makes a multi-step migration atomic and what
// guarantees the underlying connection is released on Rollback.
func (d *DB) BeginTx(ctx context.Context) (*Tx, error) {
	if d == nil || d.db == nil {
		return nil, fmt.Errorf("database unavailable")
	}
	t, e := d.db.BeginTx(ctx, nil)
	if e != nil {
		return nil, fmt.Errorf("database begin: %w", e)
	}
	return &Tx{tx: t}, nil
}
func (t *Tx) Exec(ctx context.Context, q string, args ...any) error {
	if t == nil || t.tx == nil {
		return fmt.Errorf("transaction unavailable")
	}
	if _, e := t.tx.ExecContext(ctx, q, args...); e != nil {
		return fmt.Errorf("transaction exec: %w", e)
	}
	return nil
}
func (t *Tx) Commit() error {
	if t == nil || t.tx == nil {
		return fmt.Errorf("transaction unavailable")
	}
	if e := t.tx.Commit(); e != nil {
		return fmt.Errorf("transaction commit: %w", e)
	}
	return nil
}
func (t *Tx) Rollback() error {
	if t == nil || t.tx == nil {
		return nil
	}
	if e := t.tx.Rollback(); e != nil {
		return fmt.Errorf("transaction rollback: %w", e)
	}
	return nil
}
func (d *DB) Close() error {
	if d == nil || d.db == nil {
		return nil
	}
	return d.db.Close()
}
