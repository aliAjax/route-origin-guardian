package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"testing"
)

var (
	errMigrationExec029 = errors.New("migration execution failed")
	stateMu029          sync.Mutex
	activeState029      *migrationDriverState029
)

type migrationDriverState029 struct {
	mu        sync.Mutex
	begins    int
	commits   int
	rollbacks int
	execs     []string
	cancel    context.CancelFunc
}

func (s *migrationDriverState029) snapshot() (int, int, int, []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.begins, s.commits, s.rollbacks, append([]string(nil), s.execs...)
}

type migrationDriver029 struct{}

func (migrationDriver029) Open(string) (driver.Conn, error) {
	stateMu029.Lock()
	defer stateMu029.Unlock()
	return &migrationConn029{state: activeState029}, nil
}

type migrationConn029 struct{ state *migrationDriverState029 }

func (*migrationConn029) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not supported")
}
func (*migrationConn029) Close() error { return nil }
func (c *migrationConn029) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}
func (c *migrationConn029) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	c.state.mu.Lock()
	c.state.begins++
	c.state.mu.Unlock()
	return &migrationTx029{state: c.state}, nil
}
func (c *migrationConn029) ExecContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	c.state.mu.Lock()
	c.state.execs = append(c.state.execs, query)
	cancel := c.state.cancel
	c.state.mu.Unlock()
	switch query {
	case "FAIL":
		return nil, errMigrationExec029
	case "CANCEL":
		if cancel != nil {
			cancel()
		}
		<-ctx.Done()
		return nil, ctx.Err()
	default:
		return driver.RowsAffected(1), nil
	}
}

type migrationTx029 struct{ state *migrationDriverState029 }

func (t *migrationTx029) Commit() error {
	t.state.mu.Lock()
	t.state.commits++
	t.state.mu.Unlock()
	return nil
}
func (t *migrationTx029) Rollback() error {
	t.state.mu.Lock()
	t.state.rollbacks++
	t.state.mu.Unlock()
	return nil
}

func init() { sql.Register("migration-driver-029", migrationDriver029{}) }

func newMigrationDB029(t *testing.T, state *migrationDriverState029) *DB {
	t.Helper()
	stateMu029.Lock()
	activeState029 = state
	stateMu029.Unlock()
	db, err := sql.Open("migration-driver-029", fmt.Sprintf("%p", state))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return New(db)
}

func TestMigrationTxnRejectsPlanBeforeBegin(t *testing.T) {
	state := &migrationDriverState029{}
	db := newMigrationDB029(t, state)
	err := ApplyMigrations(context.Background(), db, []Migration{{Version: 1, SQL: "CREATE"}, {Version: 0, SQL: ""}})
	if err == nil {
		t.Fatal("invalid migration plan was accepted")
	}
	begins, _, _, execs := state.snapshot()
	if begins != 0 || len(execs) != 0 {
		t.Fatalf("invalid plan touched database: begins=%d execs=%v", begins, execs)
	}
}

func TestMigrationTxnRollsBackExecFailure(t *testing.T) {
	state := &migrationDriverState029{}
	db := newMigrationDB029(t, state)
	err := ApplyMigrations(context.Background(), db, []Migration{{Version: 1, SQL: "CREATE"}, {Version: 2, SQL: "FAIL"}})
	if err == nil {
		t.Fatal("failed migration returned nil")
	}
	begins, commits, rollbacks, _ := state.snapshot()
	if begins != 1 || commits != 0 || rollbacks != 1 {
		t.Fatalf("failure lifecycle: begins=%d commits=%d rollbacks=%d", begins, commits, rollbacks)
	}
}

func TestMigrationTxnCommitsOnce(t *testing.T) {
	state := &migrationDriverState029{}
	db := newMigrationDB029(t, state)
	err := ApplyMigrations(context.Background(), db, []Migration{{Version: 1, SQL: "CREATE"}, {Version: 2, SQL: "ALTER"}})
	if err != nil {
		t.Fatal(err)
	}
	begins, commits, rollbacks, execs := state.snapshot()
	if begins != 1 || commits != 1 || rollbacks != 0 || len(execs) != 2 {
		t.Fatalf("success lifecycle: begins=%d commits=%d rollbacks=%d execs=%v", begins, commits, rollbacks, execs)
	}
}

func TestMigrationTxnPropagatesCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	state := &migrationDriverState029{cancel: cancel}
	db := newMigrationDB029(t, state)
	err := ApplyMigrations(ctx, db, []Migration{{Version: 1, SQL: "CREATE"}, {Version: 2, SQL: "CANCEL"}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ApplyMigrations error = %v, want context.Canceled", err)
	}
	_, commits, rollbacks, _ := state.snapshot()
	if commits != 0 || rollbacks != 1 {
		t.Fatalf("cancel lifecycle: commits=%d rollbacks=%d", commits, rollbacks)
	}
}

func TestMigrationTxnPreservesCause(t *testing.T) {
	state := &migrationDriverState029{}
	db := newMigrationDB029(t, state)
	err := ApplyMigrations(context.Background(), db, []Migration{{Version: 1, SQL: "FAIL"}})
	if !errors.Is(err, errMigrationExec029) {
		t.Fatalf("ApplyMigrations error = %v, want wrapped driver cause", err)
	}
}
