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

func ApplyMigrations(ctx context.Context, db *DB, migrations []Migration) error {
	for _, m := range migrations {
		if m.Version <= 0 || strings.TrimSpace(m.SQL) == "" {
			return fmt.Errorf("invalid migration %d", m.Version)
		}
		if e := db.Exec(ctx, m.SQL); e != nil {
			return fmt.Errorf("migration %d: %w", m.Version, e)
		}
	}
	return nil
}
