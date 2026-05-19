package settings

import appsettings "SkinPrice/skinprice/internal/application/settings"

type GetAppSettingsUseCase interface {
	Execute() (appsettings.AppSettings, error)
}

type SaveAppSettingsUseCase interface {
	Execute(settings appsettings.AppSettings) error
}

type Endpoints struct {
	getAppSettingsUC  GetAppSettingsUseCase
	saveAppSettingsUC SaveAppSettingsUseCase
}

func NewEndpoints(getAppSettingsUC GetAppSettingsUseCase, saveAppSettingsUC SaveAppSettingsUseCase) *Endpoints {
	return &Endpoints{
		getAppSettingsUC:  getAppSettingsUC,
		saveAppSettingsUC: saveAppSettingsUC,
	}
}

func (e *Endpoints) GetAppSettings() (AppSettingsResponse, error) {
	settings, err := e.getAppSettingsUC.Execute()
	if err != nil {
		return AppSettingsResponse{}, err
	}
	return AppSettingsResponse{
		Currency:                   settings.Currency,
		AutoRefreshIntervalSeconds: settings.AutoRefreshIntervalSeconds,
	}, nil
}

func (e *Endpoints) SaveAppSettings(payload SaveAppSettingsRequest) error {
	return e.saveAppSettingsUC.Execute(appsettings.AppSettings{
		Currency:                   payload.Currency,
		AutoRefreshIntervalSeconds: payload.AutoRefreshIntervalSeconds,
	})
}
