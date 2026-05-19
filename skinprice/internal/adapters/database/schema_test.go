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

func TestEnsureSchemaCreatesExtendedSkinsColumns(t *testing.T) {
	connection, err := New(&Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })

	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	if _, err := connection.DB().ExecContext(context.Background(), `
INSERT INTO skins (
	market_hash_name, display_name, icon_url, page_url, price_text,
	steam_page_url, steam_price_text, lisskins_page_url, lisskins_price_text, currency
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"AK-47 | Redline", "AK-47 | Redline", "icon", "steam-url", "$10.00",
		"steam-url", "$10.00", "lis-url", "$9.50", "1",
	); err != nil {
		t.Fatalf("insert into skins with extended columns: %v", err)
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

func TestEnsureSchemaMigratesLegacySkinsTable(t *testing.T) {
	connection, err := New(&Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })

	if _, err := connection.DB().ExecContext(context.Background(), `
CREATE TABLE skins (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	market_hash_name TEXT NOT NULL UNIQUE,
	display_name TEXT NOT NULL,
	icon_url TEXT NOT NULL DEFAULT '',
	page_url TEXT NOT NULL DEFAULT '',
	price_text TEXT NOT NULL DEFAULT '',
	currency TEXT NOT NULL DEFAULT '1',
	updated_at DATETIME
)`); err != nil {
		t.Fatalf("create legacy skins: %v", err)
	}

	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	if _, err := connection.DB().ExecContext(context.Background(), `
INSERT INTO skins (
	market_hash_name, display_name, icon_url, page_url, price_text,
	steam_page_url, steam_price_text, lisskins_page_url, lisskins_price_text, currency
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"AK-47 | Redline", "AK-47 | Redline", "icon", "steam-url", "$10.00",
		"steam-url", "$10.00", "lis-url", "$9.50", "1",
	); err != nil {
		t.Fatalf("insert migrated skins row: %v", err)
	}
}
