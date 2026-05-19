package skins

type UpdateSavedSkinPrice struct {
	Updater SavedSkinPriceUpdater
}

func (uc UpdateSavedSkinPrice) Execute(params UpdateSavedSkinPriceParams) (UpdateSavedSkinPriceResult, error) {
	return uc.Updater.UpdateSavedSkinPrice(params)
}

type UpdateAllSavedSkinsPrices struct {
	Updater SavedSkinPriceUpdater
}

func (uc UpdateAllSavedSkinsPrices) Execute(params UpdateAllSavedSkinsPricesParams) (UpdateAllSavedSkinsPricesResult, error) {
	return uc.Updater.UpdateAllSavedSkinsPrices(params)
}
