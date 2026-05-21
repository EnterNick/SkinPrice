package skins

import (
	"SkinPrice/skinprice/internal/application"
	"context"
)

type NewSkinsStorage interface {
	GetList(ctx context.Context, criteria SearchCriteria, params *application.Pagination) (NewSkinsList, error)
}

type SavedSkinRepository interface {
	Save(ctx context.Context, params SaveSkinParams) (SaveSkinResult, error)
	GetSavedList(ctx context.Context, params *application.Pagination) (SavedSkinsList, error)
	GetSavedSkin(ctx context.Context, marketHashName string) (SavedSkin, error)
	ListSavedSkinNames(ctx context.Context) ([]string, error)
	UpdateSavedSkinPrices(ctx context.Context, params UpdateSavedSkinPriceResult) error
	DeleteSavedSkin(ctx context.Context, params DeleteSavedSkinParams) error
}

type MarketPriceReader interface {
	GetByMarketHashName(ctx context.Context, marketHashName, currency string) (*NewSkin, error)
}

type MarketPageURLBuilder interface {
	BuildMarketPageURL(marketHashName string) string
}

type SavedSkinPriceCollector interface {
	Collect(ctx context.Context, saved SavedSkin, params UpdateSavedSkinPriceParams) (UpdateSavedSkinPriceResult, error)
}

type SavedSkinPriceUpdater interface {
	Execute(ctx context.Context, params UpdateSavedSkinPriceParams) (UpdateSavedSkinPriceResult, error)
}

type LisSkinsTokenStorage interface {
	UpsertLisSkinsToken(ctx context.Context, encrypted string) error
	GetLisSkinsToken(ctx context.Context) (string, error)
	DeleteLisSkinsToken(ctx context.Context) error
}
