package main

import (
	adapterdbskins "SkinPrice/skinprice/internal/adapters/database/skins"
	adaptersteam "SkinPrice/skinprice/internal/adapters/http/steam"
	"SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/config"
	presenterskins "SkinPrice/skinprice/internal/presenters/skins"
)

func (a *App) registerRoutes() {
	cfg := config.Load()
	storage := &adaptersteam.Storage{Client: adaptersteam.NewSteamClient(cfg), BaseURL: cfg.SteamBaseURL}
	searchNewSkinsUC := skins.SearchNewSkins{NewSkinsStorage: storage}
	saveSkinStorage := &adapterdbskins.Storage{Conn: a.backend.Factory.DBConnection()}
	saveSkinUC := skins.SaveSkin{SkinSaver: saveSkinStorage}
	a.skinsEndpoints = presenterskins.NewEndpoints(searchNewSkinsUC, saveSkinUC)
}

func (a *App) SearchNewSkins(filter presenterskins.SearchNewSkinsFilter) (presenterskins.NewSkinsResponse, error) {
	return a.skinsEndpoints.SearchNewSkins(filter)
}

func (a *App) SaveSkin(payload presenterskins.SaveSkinRequest) error {
	return a.skinsEndpoints.SaveSkin(payload)
}
