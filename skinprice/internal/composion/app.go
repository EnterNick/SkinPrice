package composion

import "SkinPrice/skinprice/internal/composion/factory"

type BackendApp struct {
	Factory *factory.Factory
}

func NewApp() (*BackendApp, error) {
	f, err := factory.NewFactory()
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
