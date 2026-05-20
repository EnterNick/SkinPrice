package settings

type AppSettingsResponse struct {
	Currency                   string `json:"currency"`
	AutoRefreshIntervalSeconds int    `json:"auto_refresh_interval_seconds"`
	SavedSkinsViewMode         string `json:"saved_skins_view_mode"`
}

type SaveAppSettingsRequest struct {
	Currency                   string `json:"currency"`
	AutoRefreshIntervalSeconds int    `json:"auto_refresh_interval_seconds"`
	SavedSkinsViewMode         string `json:"saved_skins_view_mode"`
}
