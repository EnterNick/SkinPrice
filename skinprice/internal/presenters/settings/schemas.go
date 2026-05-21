package settings

type AppSettingsResponse struct {
	Currency                   string `json:"currency"`
	AutoRefreshEnabled         bool   `json:"auto_refresh_enabled"`
	AutoRefreshIntervalSeconds int    `json:"auto_refresh_interval_seconds"`
	SavedSkinsViewMode         string `json:"saved_skins_view_mode"`
	FontFamily                 string `json:"font_family"`
	FontSizePx                 int    `json:"font_size_px"`
}

type SaveAppSettingsRequest struct {
	Currency                   string `json:"currency"`
	AutoRefreshEnabled         bool   `json:"auto_refresh_enabled"`
	AutoRefreshIntervalSeconds int    `json:"auto_refresh_interval_seconds"`
	SavedSkinsViewMode         string `json:"saved_skins_view_mode"`
	FontFamily                 string `json:"font_family"`
	FontSizePx                 int    `json:"font_size_px"`
}
