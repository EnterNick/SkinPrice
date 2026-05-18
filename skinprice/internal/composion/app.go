package composion

import (
	"SkinPrice/skinprice/internal/composion/factory"
	"log/slog"
)

type BackendApp struct {
	Factory *factory.Factory
}

func NewApp(logger *slog.Logger) (*BackendApp, error) {
	f, err := factory.NewFactory(logger)
	if err != nil {
		return nil, err
	}
	return &BackendApp{
		Factory: f,
	}, nil
}

func (app *BackendApp) Close() error {
	return app.Factory.Close()
}
