package appsettings

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/adapters/database/ent"
	entappsetting "SkinPrice/skinprice/internal/adapters/database/ent/appsetting"
	appsettings "SkinPrice/skinprice/internal/application/settings"
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"strconv"
	"time"
)

const (
	currencyKey                   = "saved_skins.currency"
	autoRefreshEnabledKey         = "saved_skins.auto_refresh_enabled"
	autoRefreshIntervalSecondsKey = "saved_skins.auto_refresh_interval_seconds"
	savedSkinsViewModeKey         = "saved_skins.view_mode"
	fontFamilyKey                 = "ui.font_family"
	fontSizePxKey                 = "ui.font_size_px"
)

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) GetAppSettings(ctx context.Context) (appsettings.AppSettings, error) {
	rows, err := s.Conn.Client().AppSetting.Query().All(ctx)
	if err != nil {
		return appsettings.AppSettings{}, errx.E("appsettings.repository.get", errx.CodeInternal, "failed to load app settings", err)
	}

	settings := appsettings.AppSettings{
		AutoRefreshEnabled: appsettings.DefaultAutoRefreshEnabled,
	}
	for _, row := range rows {
		switch row.Key {
		case currencyKey:
			settings.Currency = row.Value
		case autoRefreshEnabledKey:
			if enabled, parseErr := strconv.ParseBool(row.Value); parseErr == nil {
				settings.AutoRefreshEnabled = enabled
			}
		case autoRefreshIntervalSecondsKey:
			if interval, parseErr := strconv.Atoi(row.Value); parseErr == nil {
				settings.AutoRefreshIntervalSeconds = interval
			}
		case savedSkinsViewModeKey:
			settings.SavedSkinsViewMode = row.Value
		case fontFamilyKey:
			settings.FontFamily = row.Value
		case fontSizePxKey:
			if size, parseErr := strconv.Atoi(row.Value); parseErr == nil {
				settings.FontSizePx = size
			}
		}
	}
	return settings, nil
}

func (s *Storage) SaveAppSettings(ctx context.Context, settings appsettings.AppSettings) error {
	values := map[string]string{
		currencyKey:                   settings.Currency,
		autoRefreshEnabledKey:         strconv.FormatBool(settings.AutoRefreshEnabled),
		autoRefreshIntervalSecondsKey: strconv.Itoa(settings.AutoRefreshIntervalSeconds),
		savedSkinsViewModeKey:         settings.SavedSkinsViewMode,
		fontFamilyKey:                 settings.FontFamily,
		fontSizePxKey:                 strconv.Itoa(settings.FontSizePx),
	}
	for key, value := range values {
		if err := s.upsertValue(ctx, key, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) upsertValue(ctx context.Context, key, value string) error {
	now := time.Now().UTC()
	count, err := s.Conn.Client().AppSetting.Update().
		Where(entappsetting.Key(key)).
		SetValue(value).
		SetUpdatedAt(now).
		Save(ctx)
	if err != nil {
		return errx.E("appsettings.repository.save.update", errx.CodeInternal, "failed to update app settings", err)
	}
	if count > 0 {
		return nil
	}

	if _, err := s.Conn.Client().AppSetting.Create().
		SetKey(key).
		SetValue(value).
		SetUpdatedAt(now).
		Save(ctx); err != nil {
		if ent.IsConstraintError(err) {
			_, retryErr := s.Conn.Client().AppSetting.Update().
				Where(entappsetting.Key(key)).
				SetValue(value).
				SetUpdatedAt(now).
				Save(ctx)
			if retryErr == nil {
				return nil
			}
		}
		return errx.E("appsettings.repository.save.insert", errx.CodeInternal, "failed to save app settings", err)
	}
	return nil
}
