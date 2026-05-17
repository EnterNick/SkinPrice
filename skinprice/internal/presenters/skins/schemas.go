package skins

type SearchNewSkinsFilter struct {
	MarketHashName *string `json:"market_hash_name"`
	Limit          int     `json:"limit"`
	Offset         int     `json:"offset"`
}

type NewSkinsResponse struct {
	Items      []NewSkinItem `json:"items"`
	TotalCount int           `json:"total_count"`
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
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
	MarketHashName string `json:"market_hash_name"`
	DisplayName    string `json:"display_name"`
	IconURL        string `json:"icon_url"`
	PageURL        string `json:"page_url"`
}
