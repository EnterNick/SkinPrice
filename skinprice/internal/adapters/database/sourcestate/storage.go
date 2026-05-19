package sourcestate

import (
	"SkinPrice/skinprice/internal/adapters/database"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"database/sql"
	"time"
)

const lisSkinsSource = "lisskins"

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) UpsertLisSkinsToken(encrypted string) error {
	ctx := context.Background()
	now := time.Now()

	updateQuery := `UPDATE source_states SET api_token_encrypted = ?, updated_at = ? WHERE source = ?`
	insertQuery := `INSERT INTO source_states (source, api_token_encrypted, updated_at) VALUES (?, ?, ?)`
	if s.Conn.Dialect() == "postgres" {
		updateQuery = `UPDATE source_states SET api_token_encrypted = $1, updated_at = $2 WHERE source = $3`
		insertQuery = `INSERT INTO source_states (source, api_token_encrypted, updated_at) VALUES ($1, $2, $3)`
	}

	result, err := s.Conn.DB().ExecContext(ctx, updateQuery, encrypted, now, lisSkinsSource)
	if err != nil {
		return errx.E("sourcestate.storage.upsert_token.update", errx.CodeInternal, "failed to update lisskins token", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errx.E("sourcestate.storage.upsert_token.rows", errx.CodeInternal, "failed to inspect lisskins token update", err)
	}
	if rowsAffected > 0 {
		return nil
	}

	if _, err = s.Conn.DB().ExecContext(ctx, insertQuery, lisSkinsSource, encrypted, now); err != nil {
		if _, retryErr := s.Conn.DB().ExecContext(ctx, updateQuery, encrypted, now, lisSkinsSource); retryErr == nil {
			return nil
		}
		return errx.E("sourcestate.storage.upsert_token.insert", errx.CodeInternal, "failed to save lisskins token", err)
	}
	return nil
}

func (s *Storage) GetLisSkinsToken() (string, error) {
	ctx := context.Background()
	query := `SELECT api_token_encrypted FROM source_states WHERE source = ? ORDER BY id DESC LIMIT 1`
	if s.Conn.Dialect() == "postgres" {
		query = `SELECT api_token_encrypted FROM source_states WHERE source = $1 ORDER BY id DESC LIMIT 1`
	}

	var encrypted string
	err := s.Conn.DB().QueryRowContext(ctx, query, lisSkinsSource).Scan(&encrypted)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", appskins.ErrLisSkinsTokenMissing
		}
		return "", errx.E("sourcestate.storage.get_token", errx.CodeInternal, "failed to load lisskins token", err)
	}
	if encrypted == "" {
		return "", appskins.ErrLisSkinsTokenMissing
	}
	return encrypted, nil
}

func (s *Storage) DeleteLisSkinsToken() error {
	ctx := context.Background()
	query := `DELETE FROM source_states WHERE source = ?`
	if s.Conn.Dialect() == "postgres" {
		query = `DELETE FROM source_states WHERE source = $1`
	}
	_, err := s.Conn.DB().ExecContext(ctx, query, lisSkinsSource)
	if err != nil {
		return errx.E("sourcestate.storage.delete_token", errx.CodeInternal, "failed to delete lisskins token", err)
	}
	return nil
}
