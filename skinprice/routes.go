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
	saveSkinStorage := &adapterdbskins.Storage{Conn: a.backend.Factory.DBConnection(), SteamStorage: storage}
	saveSkinUC := skins.SaveSkin{SkinSaver: saveSkinStorage}
	getSavedSkinsUC := skins.GetSavedSkins{SavedSkinsReader: saveSkinStorage}
	updateSavedSkinPriceUC := skins.UpdateSavedSkinPrice{Updater: saveSkinStorage}
	updateAllSavedSkinsPricesUC := skins.UpdateAllSavedSkinsPrices{Updater: saveSkinStorage}
	a.skinsEndpoints = presenterskins.NewEndpoints(searchNewSkinsUC, saveSkinUC, getSavedSkinsUC, updateSavedSkinPriceUC, updateAllSavedSkinsPricesUC)
}

func (a *App) UpdateSavedSkinPrice(payload presenterskins.UpdateSavedSkinPriceRequest) error {
	return a.skinsEndpoints.UpdateSavedSkinPrice(payload)
}

func (a *App) UpdateAllSavedSkinsPrices(payload presenterskins.UpdateAllSavedSkinsPricesRequest) error {
	return a.skinsEndpoints.UpdateAllSavedSkinsPrices(payload)
}

func (a *App) SearchNewSkins(filter presenterskins.SearchNewSkinsFilter) (presenterskins.NewSkinsResponse, error) {
	return a.skinsEndpoints.SearchNewSkins(filter)
}

func (a *App) SaveSkin(payload presenterskins.SaveSkinRequest) error {
	return a.skinsEndpoints.SaveSkin(payload)
}

func (a *App) GetSavedSkins(filter presenterskins.GetSavedSkinsFilter) (presenterskins.SavedSkinsResponse, error) {
	return a.skinsEndpoints.GetSavedSkins(filter)
}
