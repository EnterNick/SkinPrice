package database

import (
	"context"
	"testing"
)

func TestEnsureSchemaCreatesSourceStatesTable(t *testing.T) {
	connection, err := New(&Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })

	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	if _, err := connection.DB().ExecContext(context.Background(), `INSERT INTO source_states (source, api_token_encrypted) VALUES (?, ?)`, "lisskins", "token"); err != nil {
		t.Fatalf("insert into source_states: %v", err)
	}
}

func TestEnsureSchemaMigratesLegacySourceStatesTable(t *testing.T) {
	connection, err := New(&Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })

	if _, err := connection.DB().ExecContext(context.Background(), `
CREATE TABLE source_states (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	api_token_encrypted TEXT NOT NULL DEFAULT '',
	updated_at DATETIME
)`); err != nil {
		t.Fatalf("create legacy source_states: %v", err)
	}

	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	if _, err := connection.DB().ExecContext(context.Background(), `INSERT INTO source_states (source, api_token_encrypted) VALUES (?, ?)`, "lisskins", "token"); err != nil {
		t.Fatalf("insert migrated source_states row: %v", err)
	}
}
