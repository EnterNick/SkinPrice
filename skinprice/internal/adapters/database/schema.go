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

func isMissingColumnIgnored(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate column name: source") ||
		strings.Contains(message, "duplicate column name: api_token_encrypted") ||
		strings.Contains(message, "duplicate column name: updated_at")
}

const sqliteSchemaQuery = `
CREATE TABLE IF NOT EXISTS skins (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	market_hash_name TEXT NOT NULL UNIQUE,
	display_name TEXT NOT NULL,
	icon_url TEXT NOT NULL DEFAULT '',
	page_url TEXT NOT NULL DEFAULT '',
	price_text TEXT NOT NULL DEFAULT '',
	currency TEXT NOT NULL DEFAULT '1',
	updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS source_states (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	source TEXT NOT NULL UNIQUE,
	api_token_encrypted TEXT NOT NULL DEFAULT '',
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
	currency TEXT NOT NULL DEFAULT '1',
	updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS source_states (
	id BIGSERIAL PRIMARY KEY,
	source TEXT NOT NULL UNIQUE,
	api_token_encrypted TEXT NOT NULL DEFAULT '',
	updated_at TIMESTAMPTZ
);
`
