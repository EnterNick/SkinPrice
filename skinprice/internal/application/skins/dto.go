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
	MarketHashName  string
	DisplayName     string
	NameColor       string
	IconURL         string
	PageURL         string
	SteamPageURL    string
	LisSkinsPageURL string
	CSTMPageURL     string
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
	NameColor      string
	IconURL        string

	SteamPageURL      string
	SteamPriceText    string
	SteamUpdatedAt    time.Time
	LisSkinsPageURL   string
	LisSkinsPriceText string
	LisSkinsUpdatedAt time.Time
	CSTMPageURL       string
	CSTMPriceText     string
	CSTMUpdatedAt     time.Time

	Prices   []PriceSnapshotView
	Currency string
}

type UpdateSavedSkinPriceParams struct {
	MarketHashName string
	Currency       string
}

type UpdateSavedSkinPriceResult struct {
	MarketHashName string
	Prices         []PriceSnapshotView

	SteamPageURL      string
	SteamPriceText    string
	SteamUpdatedAt    time.Time
	LisSkinsPageURL   string
	LisSkinsPriceText string
	LisSkinsUpdatedAt time.Time
	CSTMPageURL       string
	CSTMPriceText     string
	CSTMUpdatedAt     time.Time
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

type PriceQuote struct {
	Source      string
	SourceLabel string
	PageURL     string
	PriceText   string
	PriceCents  *int64
	Currency    string
	FetchedAt   time.Time
	Metadata    string
}

type PriceSnapshotView struct {
	Source      string
	SourceLabel string
	PageURL     string
	PriceText   string
	PriceCents  *int64
	Currency    string
	FetchedAt   time.Time
	Status      string
}

type SourceState struct {
	Source        string
	Status        string
	LastSuccessAt time.Time
	LastError     string
	LastErrorAt   time.Time
	UpdatedAt     time.Time
}

type RefreshTaskKind string

const (
	RefreshTaskManual RefreshTaskKind = "manual"
	RefreshTaskAuto   RefreshTaskKind = "auto"
)

type RefreshTask struct {
	MarketHashName string
	Currency       string
	Kind           RefreshTaskKind
}

type Diagnostics struct {
	Version      string
	DatabasePath string
	LogPath      string
	Sources      []SourceState
}
