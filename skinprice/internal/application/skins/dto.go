package skins

import "time"

type SearchCriteria struct {
	MarketHashName *string
	Source         string
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

type SaveSkinResult struct {
	Created bool
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

type UpdateSavedSkinPriceResult struct {
	MarketHashName string
	PriceText      string
	Currency       string
	UpdatedAt      time.Time
}

type UpdateAllSavedSkinsPricesParams struct {
	Currency string
}

type UpdateSavedSkinPriceFailure struct {
	MarketHashName string
	Message        string
}

type UpdateAllSavedSkinsPricesResult struct {
	UpdatedCount int
	FailedCount  int
	Failures     []UpdateSavedSkinPriceFailure
}

type DeleteSavedSkinParams struct {
	MarketHashName string
}
