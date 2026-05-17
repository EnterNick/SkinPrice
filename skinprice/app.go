package main

import (
	"SkinPrice/skinprice/internal/composion"
	presenterskins "SkinPrice/skinprice/internal/presenters/skins"
	"context"
)

type App struct {
	ctx            context.Context
	backend        *composion.BackendApp
	skinsEndpoints *presenterskins.Endpoints
}

func NewApp() *App {
	app, err := composion.NewApp()
	if err != nil {
		panic(err)
	}
	instance := &App{
		backend: app,
	}
	instance.registerRoutes()
	return instance
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown() error {
	err := a.backend.Close()
	if err != nil {
		return err
	}
	return nil
}
