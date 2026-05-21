package sourcestate

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/adapters/database/ent"
	entsourcestate "SkinPrice/skinprice/internal/adapters/database/ent/sourcestate"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"time"
)

const lisSkinsSource = "lisskins"

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) RecordSourceSuccess(ctx context.Context, source string, at time.Time) error {
	now := time.Now().UTC()
	count, err := s.Conn.Client().SourceState.Update().
		Where(entsourcestate.Source(source)).
		SetStatus("ok").
		SetLastSuccessAt(at.UTC()).
		SetLastError("").
		ClearLastErrorAt().
		SetUpdatedAt(now).
		Save(ctx)
	if err != nil {
		return errx.E("sourcestate.repository.record_success.update", errx.CodeInternal, "failed to update source state", err)
	}
	if count > 0 {
		return nil
	}
	if _, err := s.Conn.Client().SourceState.Create().
		SetSource(source).
		SetStatus("ok").
		SetLastSuccessAt(at.UTC()).
		SetLastError("").
		SetUpdatedAt(now).
		Save(ctx); err != nil {
		return errx.E("sourcestate.repository.record_success.insert", errx.CodeInternal, "failed to create source state", err)
	}
	return nil
}

func (s *Storage) RecordSourceError(ctx context.Context, source string, message string, at time.Time) error {
	now := time.Now().UTC()
	count, err := s.Conn.Client().SourceState.Update().
		Where(entsourcestate.Source(source)).
		SetStatus("error").
		SetLastError(message).
		SetLastErrorAt(at.UTC()).
		SetUpdatedAt(now).
		Save(ctx)
	if err != nil {
		return errx.E("sourcestate.repository.record_error.update", errx.CodeInternal, "failed to update source error", err)
	}
	if count > 0 {
		return nil
	}
	if _, err := s.Conn.Client().SourceState.Create().
		SetSource(source).
		SetStatus("error").
		SetLastError(message).
		SetLastErrorAt(at.UTC()).
		SetUpdatedAt(now).
		Save(ctx); err != nil {
		return errx.E("sourcestate.repository.record_error.insert", errx.CodeInternal, "failed to create source error", err)
	}
	return nil
}

func (s *Storage) ListSourceStates(ctx context.Context) ([]appskins.SourceState, error) {
	rows, err := s.Conn.Client().SourceState.Query().
		Order(ent.Asc(entsourcestate.FieldSource)).
		All(ctx)
	if err != nil {
		return nil, errx.E("sourcestate.repository.list", errx.CodeInternal, "failed to list source states", err)
	}
	result := make([]appskins.SourceState, 0, len(rows))
	for _, row := range rows {
		state := appskins.SourceState{
			Source:    row.Source,
			Status:    row.Status,
			LastError: row.LastError,
		}
		if row.LastSuccessAt != nil {
			state.LastSuccessAt = *row.LastSuccessAt
		}
		if row.LastErrorAt != nil {
			state.LastErrorAt = *row.LastErrorAt
		}
		if row.UpdatedAt != nil {
			state.UpdatedAt = *row.UpdatedAt
		}
		result = append(result, state)
	}
	return result, nil
}

func (s *Storage) UpsertLisSkinsToken(ctx context.Context, encrypted string) error {
	now := time.Now().UTC()
	count, err := s.Conn.Client().SourceState.Update().
		Where(entsourcestate.Source(lisSkinsSource)).
		SetAPITokenEncrypted(encrypted).
		SetUpdatedAt(now).
		Save(ctx)
	if err != nil {
		return errx.E("sourcestate.repository.upsert_token.update", errx.CodeInternal, "failed to update lisskins token", err)
	}
	if count > 0 {
		return nil
	}

	if _, err := s.Conn.Client().SourceState.Create().
		SetSource(lisSkinsSource).
		SetAPITokenEncrypted(encrypted).
		SetUpdatedAt(now).
		Save(ctx); err != nil {
		if ent.IsConstraintError(err) {
			_, retryErr := s.Conn.Client().SourceState.Update().
				Where(entsourcestate.Source(lisSkinsSource)).
				SetAPITokenEncrypted(encrypted).
				SetUpdatedAt(now).
				Save(ctx)
			if retryErr == nil {
				return nil
			}
		}
		return errx.E("sourcestate.repository.upsert_token.insert", errx.CodeInternal, "failed to save lisskins token", err)
	}
	return nil
}

func (s *Storage) GetLisSkinsToken(ctx context.Context) (string, error) {
	row, err := s.Conn.Client().SourceState.Query().
		Where(entsourcestate.Source(lisSkinsSource)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", appskins.ErrLisSkinsTokenMissing
		}
		return "", errx.E("sourcestate.repository.get_token", errx.CodeInternal, "failed to load lisskins token", err)
	}
	if row.APITokenEncrypted == "" {
		return "", appskins.ErrLisSkinsTokenMissing
	}
	return row.APITokenEncrypted, nil
}

func (s *Storage) DeleteLisSkinsToken(ctx context.Context) error {
	if _, err := s.Conn.Client().SourceState.Delete().
		Where(entsourcestate.Source(lisSkinsSource)).
		Exec(ctx); err != nil {
		return errx.E("sourcestate.repository.delete_token", errx.CodeInternal, "failed to delete lisskins token", err)
	}
	return nil
}
