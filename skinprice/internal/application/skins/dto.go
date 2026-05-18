package skins

import "time"

type SearchCriteria struct {
	MarketHashName *string
}

type NewSkinsList struct {
	Items      []NewSkin
	TotalCount int
	Offset     int
	Limit      int
}

type NewSkin struct {
	MarketHashName string
	DisplayName    string

	SellListings int64

	PriceCents *int64
	PriceText  string

	IconURL string
	PageURL string
}

type SaveSkinParams struct {
	MarketHashName string
	DisplayName    string
	IconURL        string
	PageURL        string
}

type SavedSkinsList struct {
	Items      []SavedSkin
	TotalCount int
	Offset     int
	Limit      int
}

type SavedSkin struct {
	MarketHashName string
	DisplayName    string
	IconURL        string
	PageURL        string
	PriceText      string
	Currency       string
	UpdatedAt      time.Time
}

type UpdateSavedSkinPriceParams struct {
	MarketHashName string
	Currency       string
}

type UpdateAllSavedSkinsPricesParams struct {
	Currency string
}
