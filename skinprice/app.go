package main

import (
	"SkinPrice/skinprice/internal/composion"
	presentersettings "SkinPrice/skinprice/internal/presenters/settings"
	presenterskins "SkinPrice/skinprice/internal/presenters/skins"
	"SkinPrice/skinprice/internal/shared/errx"
	"SkinPrice/skinprice/internal/shared/logx"
	"context"
	"log/slog"
)

type App struct {
	ctx               context.Context
	backend           *composion.BackendApp
	skinsEndpoints    *presenterskins.Endpoints
	settingsEndpoints *presentersettings.Endpoints
	logger            *slog.Logger
}

func NewApp(logger *slog.Logger) (*App, error) {
	logger = logx.WithComponent(logger, "app")
	app, err := composion.NewApp(logger)
	if err != nil {
		logger.Error("failed to initialize backend application", logx.ErrAttrs(err)...)
		return nil, errx.FromError(err, "failed to initialize application")
	}
	instance := &App{
		backend: app,
		logger:  logger,
	}
	instance.registerRoutes()
	logger.Info("application initialized")
	return instance, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logger.Info("wails startup completed")
}

func (a *App) Shutdown() error {
	err := a.backend.Close()
	if err != nil {
		a.logger.Error("application shutdown failed", logx.ErrAttrs(err)...)
		return err
	}
	a.logger.Info("application shutdown completed")
	return nil
}
