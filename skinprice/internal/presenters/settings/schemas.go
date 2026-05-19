package settings

type AppSettingsResponse struct {
	Currency                   string `json:"currency"`
	AutoRefreshIntervalSeconds int    `json:"auto_refresh_interval_seconds"`
}

type SaveAppSettingsRequest struct {
	Currency                   string `json:"currency"`
	AutoRefreshIntervalSeconds int    `json:"auto_refresh_interval_seconds"`
}
