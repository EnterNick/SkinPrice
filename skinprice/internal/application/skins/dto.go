package skins

import "time"

type SearchCriteria struct {
	MarketHashName *string
	SortColumn     string
	SortDir        string
	PriceMin       *string
	PriceMax       *string

	SearchDescriptions bool

	Types           []string
	Weapons         []string
	Rarities        []string
	Exteriors       []string
	ItemSets        []string
	ProPlayers      []string
	StickerCapsules []string
	TournamentTeams []string
}

type NewSkinsList struct {
	Items      []NewSkin
	TotalCount int
	Offset     int
	Limit      int
	NextCursor string
}

type NewSkin struct {
	MarketHashName string
	DisplayName    string
	NameColor      string

	SellListings int64

	PriceCents *int64
	PriceText  string

	IconURL string
	PageURL string
}

type SaveSkinParams struct {
	MarketHashName string
	DisplayName    string
	NameColor      string
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
	MarketHashName    string
	DisplayName       string
	NameColor         string
	IconURL           string
	SteamPageURL      string
	SteamPriceText    string
	SteamUpdatedAt    time.Time
	LisSkinsPageURL   string
	LisSkinsPriceText string
	LisSkinsUpdatedAt time.Time
	Currency          string
}

type UpdateSavedSkinPriceParams struct {
	MarketHashName string
	Currency       string
}

type UpdateSavedSkinPriceResult struct {
	MarketHashName    string
	SteamPageURL      string
	SteamPriceText    string
	SteamUpdatedAt    time.Time
	LisSkinsPageURL   string
	LisSkinsPriceText string
	LisSkinsUpdatedAt time.Time
	Currency          string
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
