package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func EnsureSchema(connection *Connection) error {
	ctx := context.Background()
	if err := ensureSchemaMigrations(ctx, connection); err != nil {
		return fmt.Errorf("ensure schema migrations: %w", err)
	}
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
	if err := ensurePriceSnapshotsSchema(ctx, connection); err != nil {
		return fmt.Errorf("ensure price_snapshots schema: %w", err)
	}
	if err := migrateLegacySkinsToPriceSnapshots(ctx, connection); err != nil {
		return fmt.Errorf("migrate legacy skins to price_snapshots: %w", err)
	}
	if err := recordSchemaMigration(ctx, connection, 1, "baseline_schema"); err != nil {
		return fmt.Errorf("record baseline schema migration: %w", err)
	}

	return nil
}

func ensureSchemaMigrations(ctx context.Context, connection *Connection) error {
	if connection.Dialect() == "postgres" {
		_, err := connection.DB().ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version BIGINT PRIMARY KEY,
	name TEXT NOT NULL DEFAULT '',
	applied_at TIMESTAMPTZ NOT NULL
)`)
		return err
	}
	_, err := connection.DB().ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version INTEGER PRIMARY KEY,
	name TEXT NOT NULL DEFAULT '',
	applied_at DATETIME NOT NULL
)`)
	return err
}

func recordSchemaMigration(ctx context.Context, connection *Connection, version int, name string) error {
	if connection.Dialect() == "postgres" {
		_, err := connection.DB().ExecContext(ctx, `
INSERT INTO schema_migrations (version, name, applied_at)
VALUES ($1, $2, CURRENT_TIMESTAMP)
ON CONFLICT (version) DO NOTHING`, version, name)
		return err
	}
	_, err := connection.DB().ExecContext(ctx, `
INSERT OR IGNORE INTO schema_migrations (version, name, applied_at)
VALUES (?, ?, CURRENT_TIMESTAMP)`, version, name)
	return err
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
		for _, statement := range []string{
			`ALTER TABLE source_states ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'unknown'`,
			`ALTER TABLE source_states ADD COLUMN IF NOT EXISTS last_success_at TIMESTAMPTZ`,
			`ALTER TABLE source_states ADD COLUMN IF NOT EXISTS last_error TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE source_states ADD COLUMN IF NOT EXISTS last_error_at TIMESTAMPTZ`,
		} {
			if _, err := connection.DB().ExecContext(ctx, statement); err != nil {
				return err
			}
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
	for _, statement := range []string{
		`ALTER TABLE source_states ADD COLUMN status TEXT NOT NULL DEFAULT 'unknown'`,
		`ALTER TABLE source_states ADD COLUMN last_success_at DATETIME`,
		`ALTER TABLE source_states ADD COLUMN last_error TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE source_states ADD COLUMN last_error_at DATETIME`,
	} {
		if _, err := connection.DB().ExecContext(ctx, statement); err != nil && !isMissingColumnIgnored(err) {
			return err
		}
	}
	if _, err := connection.DB().ExecContext(ctx, `ALTER TABLE source_states ADD COLUMN updated_at DATETIME`); err != nil && !isMissingColumnIgnored(err) {
		return err
	}
	if _, err := connection.DB().ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS source_states_source_uq ON source_states (source)`); err != nil {
		return err
	}
	return nil
}

func ensurePriceSnapshotsSchema(ctx context.Context, connection *Connection) error {
	if connection.Dialect() == "postgres" {
		statements := []string{
			`CREATE TABLE IF NOT EXISTS price_snapshots (
				id BIGSERIAL PRIMARY KEY,
				market_hash_name TEXT NOT NULL,
				source TEXT NOT NULL,
				source_label TEXT NOT NULL DEFAULT '',
				page_url TEXT NOT NULL DEFAULT '',
				price_text TEXT NOT NULL DEFAULT '',
				price_cents BIGINT,
				currency TEXT NOT NULL DEFAULT '1',
				fetched_at TIMESTAMPTZ NOT NULL,
				metadata TEXT NOT NULL DEFAULT ''
			)`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS market_hash_name TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS source_label TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS page_url TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS price_text TEXT NOT NULL DEFAULT ''`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS price_cents BIGINT`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS currency TEXT NOT NULL DEFAULT '1'`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS fetched_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`,
			`ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS metadata TEXT NOT NULL DEFAULT ''`,
			`CREATE INDEX IF NOT EXISTS price_snapshots_market_source_fetched_idx ON price_snapshots (market_hash_name, source, fetched_at)`,
			`CREATE INDEX IF NOT EXISTS price_snapshots_source_fetched_idx ON price_snapshots (source, fetched_at)`,
		}
		for _, statement := range statements {
			if _, err := connection.DB().ExecContext(ctx, statement); err != nil {
				return err
			}
		}
		return nil
	}

	statements := []string{
		`CREATE TABLE IF NOT EXISTS price_snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			market_hash_name TEXT NOT NULL,
			source TEXT NOT NULL,
			source_label TEXT NOT NULL DEFAULT '',
			page_url TEXT NOT NULL DEFAULT '',
			price_text TEXT NOT NULL DEFAULT '',
			price_cents INTEGER,
			currency TEXT NOT NULL DEFAULT '1',
			fetched_at DATETIME NOT NULL,
			metadata TEXT NOT NULL DEFAULT ''
		)`,
		`ALTER TABLE price_snapshots ADD COLUMN market_hash_name TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE price_snapshots ADD COLUMN source TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE price_snapshots ADD COLUMN source_label TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE price_snapshots ADD COLUMN page_url TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE price_snapshots ADD COLUMN price_text TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE price_snapshots ADD COLUMN price_cents INTEGER`,
		`ALTER TABLE price_snapshots ADD COLUMN currency TEXT NOT NULL DEFAULT '1'`,
		`ALTER TABLE price_snapshots ADD COLUMN fetched_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`,
		`ALTER TABLE price_snapshots ADD COLUMN metadata TEXT NOT NULL DEFAULT ''`,
		`CREATE INDEX IF NOT EXISTS price_snapshots_market_source_fetched_idx ON price_snapshots (market_hash_name, source, fetched_at)`,
		`CREATE INDEX IF NOT EXISTS price_snapshots_source_fetched_idx ON price_snapshots (source, fetched_at)`,
	}
	for _, statement := range statements {
		if _, err := connection.DB().ExecContext(ctx, statement); err != nil && !isMissingColumnIgnored(err) {
			return err
		}
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

type legacySkinRow struct {
	MarketHashName    string
	SteamPageURL      string
	SteamPriceText    string
	SteamUpdatedAt    sql.NullTime
	LisSkinsPageURL   string
	LisSkinsPriceText string
	LisSkinsUpdatedAt sql.NullTime
	CSTMPageURL       string
	CSTMPriceText     string
	CSTMUpdatedAt     sql.NullTime
	PageURL           string
	PriceText         string
	UpdatedAt         sql.NullTime
	Currency          string
}

type legacyPriceSnapshotSeed struct {
	source      string
	sourceLabel string
	pageURL     string
	priceText   string
	fetchedAt   time.Time
	currency    string
}

func migrateLegacySkinsToPriceSnapshots(ctx context.Context, connection *Connection) error {
	rows, err := connection.DB().QueryContext(ctx, `
SELECT market_hash_name,
	steam_page_url, steam_price_text, steam_updated_at,
	lisskins_page_url, lisskins_price_text, lisskins_updated_at,
	cstm_page_url, cstm_price_text, cstm_updated_at,
	page_url, price_text, updated_at, currency
FROM skins`)
	if err != nil {
		return err
	}
	defer rows.Close()

	now := time.Now().UTC()
	legacyRows := make([]legacySkinRow, 0)
	for rows.Next() {
		var row legacySkinRow
		if err := rows.Scan(
			&row.MarketHashName,
			&row.SteamPageURL, &row.SteamPriceText, &row.SteamUpdatedAt,
			&row.LisSkinsPageURL, &row.LisSkinsPriceText, &row.LisSkinsUpdatedAt,
			&row.CSTMPageURL, &row.CSTMPriceText, &row.CSTMUpdatedAt,
			&row.PageURL, &row.PriceText, &row.UpdatedAt, &row.Currency,
		); err != nil {
			return err
		}
		legacyRows = append(legacyRows, row)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_ = rows.Close()

	for _, row := range legacyRows {
		if err := migrateLegacySkinRow(ctx, connection, row, now); err != nil {
			return err
		}
	}
	return nil
}

func migrateLegacySkinRow(ctx context.Context, connection *Connection, row legacySkinRow, now time.Time) error {
	steamSeed := legacyPriceSnapshotSeed{
		source:      "steam",
		sourceLabel: "Steam",
		pageURL:     firstNonEmpty(row.SteamPageURL, row.PageURL),
		priceText:   firstNonEmpty(row.SteamPriceText, row.PriceText),
		fetchedAt:   chooseLegacyTime(row.SteamUpdatedAt, row.UpdatedAt, now),
		currency:    normalizeLegacyCurrency(row.Currency),
	}
	lisSeed := legacyPriceSnapshotSeed{
		source:      "lisskins",
		sourceLabel: "LisSkins",
		pageURL:     row.LisSkinsPageURL,
		priceText:   row.LisSkinsPriceText,
		fetchedAt:   chooseLegacyTime(row.LisSkinsUpdatedAt, row.UpdatedAt, now),
		currency:    "1",
	}
	cstmSeed := legacyPriceSnapshotSeed{
		source:      "cstm",
		sourceLabel: "CS TM",
		pageURL:     row.CSTMPageURL,
		priceText:   row.CSTMPriceText,
		fetchedAt:   chooseLegacyTime(row.CSTMUpdatedAt, row.UpdatedAt, now),
		currency:    normalizeLegacyCurrency(row.Currency),
	}

	for _, seed := range []legacyPriceSnapshotSeed{steamSeed, lisSeed, cstmSeed} {
		if seed.pageURL == "" && seed.priceText == "" {
			continue
		}
		exists, err := priceSnapshotExists(ctx, connection, row.MarketHashName, seed.source)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if err := insertLegacyPriceSnapshot(ctx, connection, row.MarketHashName, seed); err != nil {
			return err
		}
	}
	return nil
}

func priceSnapshotExists(ctx context.Context, connection *Connection, marketHashName, source string) (bool, error) {
	var exists int
	err := connection.DB().QueryRowContext(ctx, `
SELECT 1
FROM price_snapshots
WHERE market_hash_name = ? AND source = ?
LIMIT 1`, marketHashName, source).Scan(&exists)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, err
}

func insertLegacyPriceSnapshot(ctx context.Context, connection *Connection, marketHashName string, seed legacyPriceSnapshotSeed) error {
	_, err := connection.DB().ExecContext(ctx, `
INSERT INTO price_snapshots (
	market_hash_name, source, source_label, page_url, price_text, currency, fetched_at, metadata
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		marketHashName,
		seed.source,
		seed.sourceLabel,
		seed.pageURL,
		seed.priceText,
		normalizeLegacyCurrency(seed.currency),
		seed.fetchedAt.UTC(),
		"legacy-migrated",
	)
	return err
}

func chooseLegacyTime(primary, fallback sql.NullTime, defaultValue time.Time) time.Time {
	if primary.Valid {
		return primary.Time.UTC()
	}
	if fallback.Valid {
		return fallback.Time.UTC()
	}
	return defaultValue.UTC()
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func normalizeLegacyCurrency(value string) string {
	switch strings.TrimSpace(value) {
	case "3", "EUR":
		return "3"
	case "5", "RUB":
		return "5"
	default:
		return "1"
	}
}

func isMissingColumnIgnored(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate column name: source") ||
		strings.Contains(message, "duplicate column name: api_token_encrypted") ||
		strings.Contains(message, "duplicate column name: status") ||
		strings.Contains(message, "duplicate column name: last_success_at") ||
		strings.Contains(message, "duplicate column name: last_error") ||
		strings.Contains(message, "duplicate column name: last_error_at") ||
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
		strings.Contains(message, "duplicate column name: cstm_updated_at") ||
		strings.Contains(message, "duplicate column name: market_hash_name") ||
		strings.Contains(message, "duplicate column name: source_label") ||
		strings.Contains(message, "duplicate column name: page_url") ||
		strings.Contains(message, "duplicate column name: price_text") ||
		strings.Contains(message, "duplicate column name: price_cents") ||
		strings.Contains(message, "duplicate column name: currency") ||
		strings.Contains(message, "duplicate column name: fetched_at") ||
		strings.Contains(message, "duplicate column name: metadata")
}
