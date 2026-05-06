package main

import (
	"SkinPrice/skinprice/internal/composion"
	"context"
)

type App struct {
	ctx     context.Context
	backend *composion.BackendApp
}

func NewApp() *App {
	app, err := composion.NewApp()
	if err != nil {
		panic(err)
	}
	return &App{
		backend: app,
	}
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
