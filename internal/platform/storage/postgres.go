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
		return fmt.Errorf("database exec: %v", e)
	}
	return nil
}
func (d *DB) Close() error {
	if d == nil || d.db == nil {
		return nil
	}
	return d.db.Close()
}
