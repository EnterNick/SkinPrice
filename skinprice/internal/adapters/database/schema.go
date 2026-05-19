package database

import (
	"context"
	"fmt"
)

func EnsureSchema(connection *Connection) error {
	query := sqliteSchemaQuery
	if connection.Dialect() == "postgres" {
		query = postgresSchemaQuery
	}

	if _, err := connection.DB().ExecContext(context.Background(), query); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	return nil
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
