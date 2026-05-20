package skins

import "time"

type SearchNewSkinsFilter struct {
	MarketHashName     *string  `json:"market_hash_name"`
	SortColumn         string   `json:"sort_column"`
	SortDir            string   `json:"sort_dir"`
	PriceMin           *string  `json:"price_min"`
	PriceMax           *string  `json:"price_max"`
	SearchDescriptions bool     `json:"search_descriptions"`
	Type               []string `json:"type"`
	Weapon             []string `json:"weapon"`
	Rarity             []string `json:"rarity"`
	Exterior           []string `json:"exterior"`
	ItemSet            []string `json:"item_set"`
	ProPlayer          []string `json:"pro_player"`
	StickerCapsule     []string `json:"sticker_capsule"`
	TournamentTeam     []string `json:"tournament_team"`
	Limit              int      `json:"limit"`
	Offset             int      `json:"offset"`
	Cursor             string   `json:"cursor"`
}

type NewSkinsResponse struct {
	Items      []NewSkinItem `json:"items"`
	TotalCount int           `json:"total_count"`
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
	NextCursor string        `json:"next_cursor"`
}

type NewSkinItem struct {
	MarketHashName string `json:"market_hash_name"`
	DisplayName    string `json:"display_name"`
	NameColor      string `json:"name_color"`
	SellListings   int64  `json:"sell_listings"`
	PriceCents     *int64 `json:"price_cents"`
	PriceText      string `json:"price_text"`
	IconURL        string `json:"icon_url"`
	PageURL        string `json:"page_url"`
}

type SaveSkinRequest struct {
	MarketHashName string `json:"market_hash_name"`
	DisplayName    string `json:"display_name"`
	NameColor      string `json:"name_color"`
	IconURL        string `json:"icon_url"`
	PageURL        string `json:"page_url"`
}

type SaveSkinResponse struct {
	Created bool `json:"created"`
}

type GetSavedSkinsFilter struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SavedSkinsResponse struct {
	Items      []SavedSkinItem `json:"items"`
	TotalCount int             `json:"total_count"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
}

type SavedSkinItem struct {
	MarketHashName    string    `json:"market_hash_name"`
	DisplayName       string    `json:"display_name"`
	NameColor         string    `json:"name_color"`
	IconURL           string    `json:"icon_url"`
	SteamPageURL      string    `json:"steam_page_url"`
	SteamPriceText    string    `json:"steam_price_text"`
	SteamUpdatedAt    time.Time `json:"steam_updated_at"`
	LisSkinsPageURL   string    `json:"lisskins_page_url"`
	LisSkinsPriceText string    `json:"lisskins_price_text"`
	LisSkinsUpdatedAt time.Time `json:"lisskins_updated_at"`
	Currency          string    `json:"currency"`
}

type UpdateSavedSkinPriceRequest struct {
	MarketHashName string `json:"market_hash_name"`
	Currency       string `json:"currency"`
}

type UpdateSavedSkinPriceResponse struct {
	MarketHashName    string    `json:"market_hash_name"`
	SteamPageURL      string    `json:"steam_page_url"`
	SteamPriceText    string    `json:"steam_price_text"`
	SteamUpdatedAt    time.Time `json:"steam_updated_at"`
	LisSkinsPageURL   string    `json:"lisskins_page_url"`
	LisSkinsPriceText string    `json:"lisskins_price_text"`
	LisSkinsUpdatedAt time.Time `json:"lisskins_updated_at"`
	Currency          string    `json:"currency"`
}

type UpdateAllSavedSkinsPricesRequest struct {
	Currency string `json:"currency"`
}

type UpdateSavedSkinPriceFailure struct {
	MarketHashName string `json:"market_hash_name"`
	Message        string `json:"message"`
}

type UpdateAllSavedSkinsPricesResponse struct {
	UpdatedCount int                           `json:"updated_count"`
	FailedCount  int                           `json:"failed_count"`
	Failures     []UpdateSavedSkinPriceFailure `json:"failures"`
}

type DeleteSavedSkinRequest struct {
	MarketHashName string `json:"market_hash_name"`
}

type SetLisSkinsTokenRequest struct {
	Token string `json:"token"`
}

type LisSkinsTokenStatusResponse struct {
	HasToken bool `json:"hasToken"`
}
