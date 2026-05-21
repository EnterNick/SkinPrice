package database

import (
	"context"
	"fmt"
	"strings"
)

func EnsureSchema(connection *Connection) error {
	ctx := context.Background()
	if err := connection.Client().Schema.Create(ctx); err != nil {
		return fmt.Errorf("apply ent schema: %w", err)
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
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS name_color TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS steam_page_url TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS steam_price_text TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS steam_updated_at TIMESTAMPTZ`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS lisskins_page_url TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS lisskins_price_text TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS lisskins_updated_at TIMESTAMPTZ`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS cstm_page_url TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS cstm_price_text TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE skins ADD COLUMN IF NOT EXISTS cstm_updated_at TIMESTAMPTZ`,
		}
		for _, statement := range statements {
			if _, err := connection.DB().ExecContext(ctx, statement); err != nil {
				return err
			}
		}
		return nil
	}

	statements := []string{
		`ALTER TABLE skins ADD COLUMN name_color TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN steam_page_url TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN steam_price_text TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN steam_updated_at DATETIME`,
		`ALTER TABLE skins ADD COLUMN lisskins_page_url TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN lisskins_price_text TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN lisskins_updated_at DATETIME`,
		`ALTER TABLE skins ADD COLUMN cstm_page_url TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN cstm_price_text TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE skins ADD COLUMN cstm_updated_at DATETIME`,
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
		strings.Contains(message, "duplicate column name: name_color") ||
		strings.Contains(message, "duplicate column name: steam_page_url") ||
		strings.Contains(message, "duplicate column name: steam_price_text") ||
		strings.Contains(message, "duplicate column name: steam_updated_at") ||
		strings.Contains(message, "duplicate column name: lisskins_page_url") ||
		strings.Contains(message, "duplicate column name: lisskins_price_text") ||
		strings.Contains(message, "duplicate column name: lisskins_updated_at") ||
		strings.Contains(message, "duplicate column name: cstm_page_url") ||
		strings.Contains(message, "duplicate column name: cstm_price_text") ||
		strings.Contains(message, "duplicate column name: cstm_updated_at")
}
