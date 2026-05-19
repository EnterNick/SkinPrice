package skins

import "time"

type SearchNewSkinsFilter struct {
	MarketHashName *string `json:"market_hash_name"`
	Limit          int     `json:"limit"`
	Offset         int     `json:"offset"`
	Cursor         string  `json:"cursor"`
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
	SellListings   int64  `json:"sell_listings"`
	PriceCents     *int64 `json:"price_cents"`
	PriceText      string `json:"price_text"`
	IconURL        string `json:"icon_url"`
	PageURL        string `json:"page_url"`
}

type SaveSkinRequest struct {
	MarketHashName string `json:"market_hash_name"`
	DisplayName    string `json:"display_name"`
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
