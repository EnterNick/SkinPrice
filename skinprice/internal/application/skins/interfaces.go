package skins

import "SkinPrice/skinprice/internal/application"

type NewSkinsStorage interface {
	GetList(criteria SearchCriteria, params *application.Pagination) (NewSkinsList, error)
}

type SkinSaver interface {
	Save(params SaveSkinParams) (SaveSkinResult, error)
}

type SavedSkinsReader interface {
	GetSavedList(params *application.Pagination) (SavedSkinsList, error)
}

type SavedSkinPriceUpdater interface {
	UpdateSavedSkinPrice(params UpdateSavedSkinPriceParams) (UpdateSavedSkinPriceResult, error)
	UpdateAllSavedSkinsPrices(params UpdateAllSavedSkinsPricesParams) (UpdateAllSavedSkinsPricesResult, error)
}

type SavedSkinDeleter interface {
	DeleteSavedSkin(params DeleteSavedSkinParams) error
}

type LisSkinsTokenStorage interface {
	UpsertLisSkinsToken(encrypted string) error
	GetLisSkinsToken() (string, error)
	DeleteLisSkinsToken() error
}
