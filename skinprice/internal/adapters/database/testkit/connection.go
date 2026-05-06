package testkit

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/adapters/database/ent"
	"context"
	"sync"
	"testing"
)

var (
	once    sync.Once
	conn    *database.Connection
	initErr error
)

func BaseConn(t *testing.T) *database.Connection {
	t.Helper()
	once.Do(func() {
		cfg := loadConfig()
		conn, initErr = database.New(cfg)
	})
	if initErr != nil {
		t.Fatalf("db init: %v", initErr)
	}
	return conn
}

func BaseClient(t *testing.T) *ent.Client {
	t.Helper()
	return BaseConn(t).Client()
}

func Close() error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func WithTx(t *testing.T, fn func(ctx context.Context, c *ent.Client)) {
	t.Helper()
	ctx := t.Context()

	tx, err := BaseClient(t).Tx(ctx)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })

	fn(ctx, tx.Client())
}
