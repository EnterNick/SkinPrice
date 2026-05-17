package skins

type UpdateSavedSkinPrice struct {
	Updater SavedSkinPriceUpdater
}

func (uc UpdateSavedSkinPrice) Execute(params UpdateSavedSkinPriceParams) error {
	return uc.Updater.UpdateSavedSkinPrice(params)
}

type UpdateAllSavedSkinsPrices struct {
	Updater SavedSkinPriceUpdater
}

func (uc UpdateAllSavedSkinsPrices) Execute(params UpdateAllSavedSkinsPricesParams) error {
	return uc.Updater.UpdateAllSavedSkinsPrices(params)
}
