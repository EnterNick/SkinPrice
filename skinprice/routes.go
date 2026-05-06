package main

import (
	adaptersteam "SkinPrice/skinprice/internal/adapters/http/steam"
	"SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/config"
	presenterskins "SkinPrice/skinprice/internal/presenters/skins"
)

func (a *App) registerRoutes() {
	cfg := config.Load()
	storage := &adaptersteam.Storage{Client: adaptersteam.NewSteamClient(cfg), BaseURL: cfg.SteamBaseURL}
	searchNewSkinsUC := skins.SearchNewSkins{NewSkinsStorage: storage}
	a.skinsEndpoints = presenterskins.NewEndpoints(searchNewSkinsUC)
}

func (a *App) SearchNewSkins(filter presenterskins.SearchNewSkinsFilter) (presenterskins.NewSkinsResponse, error) {
	return a.skinsEndpoints.SearchNewSkins(filter)
}
