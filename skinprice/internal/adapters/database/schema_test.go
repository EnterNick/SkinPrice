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

func TestEnsureSchemaCreatesAppSettingsTable(t *testing.T) {
	connection, err := New(&Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })

	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	if _, err := connection.DB().ExecContext(context.Background(), `INSERT INTO app_settings (key, value) VALUES (?, ?)`, "saved_skins.currency", "1"); err != nil {
		t.Fatalf("insert into app_settings: %v", err)
	}
	if _, err := connection.DB().ExecContext(context.Background(), `INSERT INTO app_settings (key, value) VALUES (?, ?)`, "saved_skins.view_mode", "cards"); err != nil {
		t.Fatalf("insert view mode into app_settings: %v", err)
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
	market_hash_name, display_name, name_color, icon_url, page_url, price_text,
	steam_page_url, steam_price_text, lisskins_page_url, lisskins_price_text, currency
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"AK-47 | Redline", "AK-47 | Redline", "8847ff", "icon", "steam-url", "$10.00",
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
	market_hash_name, display_name, name_color, icon_url, page_url, price_text,
	steam_page_url, steam_price_text, lisskins_page_url, lisskins_price_text, currency
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"AK-47 | Redline", "AK-47 | Redline", "8847ff", "icon", "steam-url", "$10.00",
		"steam-url", "$10.00", "lis-url", "$9.50", "1",
	); err != nil {
		t.Fatalf("insert migrated skins row: %v", err)
	}
}

func TestEnsureSchemaMigratesLegacyPriceSnapshotsTable(t *testing.T) {
	connection, err := New(&Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })

	if _, err := connection.DB().ExecContext(context.Background(), `
CREATE TABLE price_snapshots (
	id INTEGER PRIMARY KEY AUTOINCREMENT
)`); err != nil {
		t.Fatalf("create legacy price_snapshots: %v", err)
	}

	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}
	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema should be idempotent: %v", err)
	}

	if _, err := connection.DB().ExecContext(context.Background(), `
INSERT INTO price_snapshots (
	market_hash_name, source, source_label, page_url, price_text, currency, fetched_at, metadata
) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?)`,
		"AK-47 | Redline", "steam", "Steam", "steam-url", "$10.00", "1", "{}",
	); err != nil {
		t.Fatalf("insert migrated price snapshot: %v", err)
	}
}

func TestEnsureSchemaRecordsMigrationWithoutTouchingLegacyRows(t *testing.T) {
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
	if _, err := connection.DB().ExecContext(context.Background(), `
INSERT INTO skins (market_hash_name, display_name, icon_url, page_url, price_text, currency)
VALUES (?, ?, ?, ?, ?, ?)`,
		"AK-47 | Redline", "AK-47 | Redline", "icon", "steam-url", "$10.00", "1",
	); err != nil {
		t.Fatalf("insert legacy skin: %v", err)
	}

	if err := EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	var count int
	if err := connection.DB().QueryRowContext(context.Background(), `SELECT COUNT(*) FROM schema_migrations WHERE version = 1`).Scan(&count); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected migration version 1 to be recorded once, got %d", count)
	}

	var snapshotCount int
	if err := connection.DB().QueryRowContext(context.Background(), `SELECT COUNT(*) FROM price_snapshots WHERE market_hash_name = ?`, "AK-47 | Redline").Scan(&snapshotCount); err != nil {
		t.Fatalf("query price_snapshots: %v", err)
	}
	if snapshotCount != 1 {
		t.Fatalf("expected one migrated snapshot, got %d", snapshotCount)
	}

	var source, pageURL, priceText string
	if err := connection.DB().QueryRowContext(context.Background(), `SELECT source, page_url, price_text FROM price_snapshots WHERE market_hash_name = ?`, "AK-47 | Redline").Scan(&source, &pageURL, &priceText); err != nil {
		t.Fatalf("query migrated price snapshot: %v", err)
	}
	if source != "steam" {
		t.Fatalf("expected migrated snapshot source steam, got %q", source)
	}
	if pageURL != "steam-url" {
		t.Fatalf("expected migrated snapshot page url, got %q", pageURL)
	}
	if priceText != "$10.00" {
		t.Fatalf("expected migrated snapshot price, got %q", priceText)
	}

	if err := connection.DB().QueryRowContext(context.Background(), `SELECT price_text FROM skins WHERE market_hash_name = ?`, "AK-47 | Redline").Scan(&priceText); err != nil {
		t.Fatalf("query legacy skin: %v", err)
	}
	if priceText != "$10.00" {
		t.Fatalf("expected legacy skin price to be preserved, got %q", priceText)
	}
}
