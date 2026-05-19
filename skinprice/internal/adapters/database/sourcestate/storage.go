package sourcestate

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/adapters/database/ent"
	"SkinPrice/skinprice/internal/adapters/database/ent/sourcestate"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"time"
)

const lisSkinsSource = "lisskins"

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) UpsertLisSkinsToken(encrypted string) error {
	ctx := context.Background()
	state, err := s.Conn.Client().SourceState.Query().Where(sourcestate.SourceEQ(lisSkinsSource)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return errx.E("sourcestate.storage.upsert_token.get", errx.CodeInternal, "failed to read lisskins token state", err)
		}
		if createErr := s.Conn.Client().SourceState.Create().
			SetSource(lisSkinsSource).
			SetAPITokenEncrypted(encrypted).
			SetUpdatedAt(time.Now()).
			Exec(ctx); createErr != nil {
			return errx.E("sourcestate.storage.upsert_token.create", errx.CodeInternal, "failed to create lisskins token", createErr)
		}
		return nil
	}

	if updateErr := s.Conn.Client().SourceState.UpdateOneID(state.ID).
		SetAPITokenEncrypted(encrypted).
		SetUpdatedAt(time.Now()).
		Exec(ctx); updateErr != nil {
		return errx.E("sourcestate.storage.upsert_token.update", errx.CodeInternal, "failed to update lisskins token", updateErr)
	}
	return nil
}

func (s *Storage) GetLisSkinsToken() (string, error) {
	ctx := context.Background()
	state, err := s.Conn.Client().SourceState.Query().Where(sourcestate.SourceEQ(lisSkinsSource)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", appskins.ErrLisSkinsTokenMissing
		}
		return "", errx.E("sourcestate.storage.get_token", errx.CodeInternal, "failed to load lisskins token", err)
	}
	if state.APITokenEncrypted == "" {
		return "", appskins.ErrLisSkinsTokenMissing
	}
	return state.APITokenEncrypted, nil
}

func (s *Storage) DeleteLisSkinsToken() error {
	ctx := context.Background()
	_, err := s.Conn.Client().SourceState.Delete().Where(sourcestate.SourceEQ(lisSkinsSource)).Exec(ctx)
	if err != nil {
		return errx.E("sourcestate.storage.delete_token", errx.CodeInternal, "failed to delete lisskins token", err)
	}
	return nil
}
