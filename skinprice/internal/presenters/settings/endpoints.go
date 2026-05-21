package settings

import (
	appsettings "SkinPrice/skinprice/internal/application/settings"
	"context"
)

type GetAppSettingsUseCase interface {
	Execute(ctx context.Context) (appsettings.AppSettings, error)
}

type SaveAppSettingsUseCase interface {
	Execute(ctx context.Context, settings appsettings.AppSettings) error
}

type Endpoints struct {
	getAppSettingsUC  GetAppSettingsUseCase
	saveAppSettingsUC SaveAppSettingsUseCase
}

type EndpointDeps struct {
	GetAppSettings  GetAppSettingsUseCase
	SaveAppSettings SaveAppSettingsUseCase
}

func NewEndpoints(deps EndpointDeps) *Endpoints {
	return &Endpoints{
		getAppSettingsUC:  deps.GetAppSettings,
		saveAppSettingsUC: deps.SaveAppSettings,
	}
}

func (e *Endpoints) GetAppSettings(ctx context.Context) (AppSettingsResponse, error) {
	settings, err := e.getAppSettingsUC.Execute(ctx)
	if err != nil {
		return AppSettingsResponse{}, err
	}
	return AppSettingsResponse{
		Currency:                   settings.Currency,
		AutoRefreshEnabled:         settings.AutoRefreshEnabled,
		AutoRefreshIntervalSeconds: settings.AutoRefreshIntervalSeconds,
		SavedSkinsViewMode:         settings.SavedSkinsViewMode,
		FontFamily:                 settings.FontFamily,
		FontSizePx:                 settings.FontSizePx,
	}, nil
}

func (e *Endpoints) SaveAppSettings(ctx context.Context, payload SaveAppSettingsRequest) error {
	return e.saveAppSettingsUC.Execute(ctx, appsettings.AppSettings{
		Currency:                   payload.Currency,
		AutoRefreshEnabled:         payload.AutoRefreshEnabled,
		AutoRefreshIntervalSeconds: payload.AutoRefreshIntervalSeconds,
		SavedSkinsViewMode:         payload.SavedSkinsViewMode,
		FontFamily:                 payload.FontFamily,
		FontSizePx:                 payload.FontSizePx,
	})
}
