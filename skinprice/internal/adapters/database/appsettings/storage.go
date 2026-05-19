package appsettings

import (
	"SkinPrice/skinprice/internal/adapters/database"
	appsettings "SkinPrice/skinprice/internal/application/settings"
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"strconv"
	"time"
)

const (
	currencyKey                   = "saved_skins.currency"
	autoRefreshIntervalSecondsKey = "saved_skins.auto_refresh_interval_seconds"
)

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) GetAppSettings() (appsettings.AppSettings, error) {
	ctx := context.Background()
	query := `SELECT key, value FROM app_settings`

	rows, err := s.Conn.DB().QueryContext(ctx, query)
	if err != nil {
		return appsettings.AppSettings{}, errx.E("appsettings.storage.get.query", errx.CodeInternal, "failed to load app settings", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	settings := appsettings.AppSettings{}
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			return appsettings.AppSettings{}, errx.E("appsettings.storage.get.scan", errx.CodeInternal, "failed to read app settings", err)
		}
		switch key {
		case currencyKey:
			settings.Currency = value
		case autoRefreshIntervalSecondsKey:
			interval, parseErr := strconv.Atoi(value)
			if parseErr == nil {
				settings.AutoRefreshIntervalSeconds = interval
			}
		}
	}
	if err := rows.Err(); err != nil {
		return appsettings.AppSettings{}, errx.E("appsettings.storage.get.rows", errx.CodeInternal, "failed to iterate app settings", err)
	}

	return settings, nil
}

func (s *Storage) SaveAppSettings(settings appsettings.AppSettings) error {
	if err := s.upsertValue(currencyKey, settings.Currency); err != nil {
		return err
	}
	if err := s.upsertValue(autoRefreshIntervalSecondsKey, strconv.Itoa(settings.AutoRefreshIntervalSeconds)); err != nil {
		return err
	}
	return nil
}

func (s *Storage) upsertValue(key, value string) error {
	ctx := context.Background()
	now := time.Now().UTC()

	updateQuery := `UPDATE app_settings SET value = ?, updated_at = ? WHERE key = ?`
	insertQuery := `INSERT INTO app_settings (key, value, updated_at) VALUES (?, ?, ?)`
	if s.Conn.Dialect() == "postgres" {
		updateQuery = `UPDATE app_settings SET value = $1, updated_at = $2 WHERE key = $3`
		insertQuery = `INSERT INTO app_settings (key, value, updated_at) VALUES ($1, $2, $3)`
	}

	result, err := s.Conn.DB().ExecContext(ctx, updateQuery, value, now, key)
	if err != nil {
		return errx.E("appsettings.storage.save.update", errx.CodeInternal, "failed to update app settings", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errx.E("appsettings.storage.save.rows", errx.CodeInternal, "failed to inspect app settings update", err)
	}
	if rowsAffected > 0 {
		return nil
	}

	if _, err = s.Conn.DB().ExecContext(ctx, insertQuery, key, value, now); err != nil {
		if _, retryErr := s.Conn.DB().ExecContext(ctx, updateQuery, value, now, key); retryErr == nil {
			return nil
		}
		return errx.E("appsettings.storage.save.insert", errx.CodeInternal, "failed to save app settings", err)
	}
	return nil
}
