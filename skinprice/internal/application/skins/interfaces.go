package skins

import "SkinPrice/skinprice/internal/application"

type NewSkinsStorage interface {
	GetList(criteria SearchCriteria, params *application.Pagination) (NewSkinsList, error)
}

type SkinSaver interface {
	Save(params SaveSkinParams) error
}

type SavedSkinsReader interface {
	GetSavedList(params *application.Pagination) (SavedSkinsList, error)
}

type SavedSkinPriceUpdater interface {
	UpdateSavedSkinPrice(params UpdateSavedSkinPriceParams) error
	UpdateAllSavedSkinsPrices(params UpdateAllSavedSkinsPricesParams) error
}
