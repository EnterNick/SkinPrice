package composion

import (
	"SkinPrice/skinprice/internal/composion/factory"
	presentersettings "SkinPrice/skinprice/internal/presenters/settings"
	presenterskins "SkinPrice/skinprice/internal/presenters/skins"
	"log/slog"
)

type BackendApp struct {
	Factory           *factory.Factory
	SkinsEndpoints    *presenterskins.Endpoints
	SettingsEndpoints *presentersettings.Endpoints
}

func NewApp(logger *slog.Logger) (*BackendApp, error) {
	f, err := factory.NewFactory(logger)
	if err != nil {
		return nil, err
	}
	return &BackendApp{
		Factory:           f,
		SkinsEndpoints:    f.SkinsEndpoints(),
		SettingsEndpoints: f.SettingsEndpoints(),
	}, nil
}

func (app *BackendApp) Close() error {
	return app.Factory.Close()
}
