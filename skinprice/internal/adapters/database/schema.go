package database

import (
	"context"
	"fmt"
	"strings"
)

func EnsureSchema(connection *Connection) error {
	ctx := context.Background()
	query := sqliteSchemaQuery
	if connection.Dialect() == "postgres" {
		query = postgresSchemaQuery
	}

	if _, err := connection.DB().ExecContext(ctx, query); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	if err := ensureSourceStatesSchema(ctx, connection); err != nil {
		return fmt.Errorf("ensure source_states schema: %w", err)
	}
	if err := ensureSkinsSchema(ctx, connection); err != nil {
		return fmt.Errorf("ensure skins schema: %w", err)
	}
	if err := ensureAppSettingsSchema(ctx, connection); err != nil {
		return fmt.Errorf("ensure app_settings schema: %w", err)
	}

	return nil
}

func ensureSkinsSchema(ctx context.Context, connection *Connection) error {
	if connection.Dialect() == "postgres" {
		statements := []string{
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS steam_page_url TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS steam_price_text TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS steam_updated_at TIMESTAMPTZ`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS lisskins_page_url TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS lisskins_price_text TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS lisskins_updated_at TIMESTAMPTZ`,
		}
		for _, statement := range statements {
			if _, err := connection.DB().ExecContext(ctx, statement); err != nil {
				return err
			}
		}
		return nil
	}

	statements := []string{
		`ALTER TABLE skins ADD COLUMN steam_page_url TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN steam_price_text TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN steam_updated_at DATETIME`,
		`ALTER TABLE skins ADD COLUMN lisskins_page_url TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN lisskins_price_text TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN lisskins_updated_at DATETIME`,
	}
	for _, statement := range statements {
		if _, err := connection.DB().ExecContext(ctx, statement); err != nil && !isMissingColumnIgnored(err) {
			return err
		}
	}
	return nil
}

func ensureSourceStatesSchema(ctx context.Context, connection *Connection) error {
	if connection.Dialect() == "postgres" {
		if _, err := connection.DB().ExecContext(ctx, `
ALTER TABLE source_states
ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'lisskins'`); err != nil {
			return err
		}
		if _, err := connection.DB().ExecContext(ctx, `
ALTER TABLE source_states
ADD COLUMN IF NOT EXISTS api_token_encrypted TEXT NOT NULL DEFAULT ''`); err != nil {
			return err
		}
		if _, err := connection.DB().ExecContext(ctx, `
ALTER TABLE source_states
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ`); err != nil {
			return err
		}
		if _, err := connection.DB().ExecContext(ctx, `
CREATE UNIQUE INDEX IF NOT EXISTS source_states_source_uq ON source_states (source)`); err != nil {
			return err
		}
		return nil
	}

	if _, err := connection.DB().ExecContext(ctx, `ALTER TABLE source_states ADD COLUMN source TEXT NOT NULL DEFAULT 'lisskins'`); err != nil && !isMissingColumnIgnored(err) {
		return err
	}
	if _, err := connection.DB().ExecContext(ctx, `ALTER TABLE source_states ADD COLUMN api_token_encrypted TEXT NOT NULL DEFAULT ''`); err != nil && !isMissingColumnIgnored(err) {
		return err
	}
	if _, err := connection.DB().ExecContext(ctx, `ALTER TABLE source_states ADD COLUMN updated_at DATETIME`); err != nil && !isMissingColumnIgnored(err) {
		return err
	}
	if _, err := connection.DB().ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS source_states_source_uq ON source_states (source)`); err != nil {
		return err
	}
	return nil
}

func ensureAppSettingsSchema(ctx context.Context, connection *Connection) error {
	if connection.Dialect() == "postgres" {
		if _, err := connection.DB().ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS app_settings (
	id BIGSERIAL PRIMARY KEY,
	key TEXT NOT NULL UNIQUE,
	value TEXT NOT NULL DEFAULT '',
	updated_at TIMESTAMPTZ
)`); err != nil {
			return err
		}
		return nil
	}

	if _, err := connection.DB().ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS app_settings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	key TEXT NOT NULL UNIQUE,
	value TEXT NOT NULL DEFAULT '',
	updated_at DATETIME
)`); err != nil {
		return err
	}
	return nil
}

func isMissingColumnIgnored(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate column name: source") ||
		strings.Contains(message, "duplicate column name: api_token_encrypted") ||
		strings.Contains(message, "duplicate column name: updated_at") ||
		strings.Contains(message, "duplicate column name: steam_page_url") ||
		strings.Contains(message, "duplicate column name: steam_price_text") ||
		strings.Contains(message, "duplicate column name: steam_updated_at") ||
		strings.Contains(message, "duplicate column name: lisskins_page_url") ||
		strings.Contains(message, "duplicate column name: lisskins_price_text") ||
		strings.Contains(message, "duplicate column name: lisskins_updated_at")
}

const sqliteSchemaQuery = `
CREATE TABLE IF NOT EXISTS skins (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	market_hash_name TEXT NOT NULL UNIQUE,
	display_name TEXT NOT NULL,
	icon_url TEXT NOT NULL DEFAULT '',
	page_url TEXT NOT NULL DEFAULT '',
	price_text TEXT NOT NULL DEFAULT '',
	steam_page_url TEXT NOT NULL DEFAULT '',
	steam_price_text TEXT NOT NULL DEFAULT '',
	steam_updated_at DATETIME,
	lisskins_page_url TEXT NOT NULL DEFAULT '',
	lisskins_price_text TEXT NOT NULL DEFAULT '',
	lisskins_updated_at DATETIME,
	currency TEXT NOT NULL DEFAULT '1',
	updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS source_states (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	source TEXT NOT NULL UNIQUE,
	api_token_encrypted TEXT NOT NULL DEFAULT '',
	updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS app_settings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	key TEXT NOT NULL UNIQUE,
	value TEXT NOT NULL DEFAULT '',
	updated_at DATETIME
);
`

const postgresSchemaQuery = `
CREATE TABLE IF NOT EXISTS skins (
	id BIGSERIAL PRIMARY KEY,
	market_hash_name TEXT NOT NULL UNIQUE,
	display_name TEXT NOT NULL,
	icon_url TEXT NOT NULL DEFAULT '',
	page_url TEXT NOT NULL DEFAULT '',
	price_text TEXT NOT NULL DEFAULT '',
	steam_page_url TEXT NOT NULL DEFAULT '',
	steam_price_text TEXT NOT NULL DEFAULT '',
	steam_updated_at TIMESTAMPTZ,
	lisskins_page_url TEXT NOT NULL DEFAULT '',
	lisskins_price_text TEXT NOT NULL DEFAULT '',
	lisskins_updated_at TIMESTAMPTZ,
	currency TEXT NOT NULL DEFAULT '1',
	updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS source_states (
	id BIGSERIAL PRIMARY KEY,
	source TEXT NOT NULL UNIQUE,
	api_token_encrypted TEXT NOT NULL DEFAULT '',
	updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS app_settings (
	id BIGSERIAL PRIMARY KEY,
	key TEXT NOT NULL UNIQUE,
	value TEXT NOT NULL DEFAULT '',
	updated_at TIMESTAMPTZ
);
`
