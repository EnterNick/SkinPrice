package steam

import (
	"SkinPrice/skinprice/internal/config"
	"net/http"
)

func NewSteamClient(config config.Config) *http.Client {
	tr := &http.Transport{
		IdleConnTimeout: config.HTTPTimeout,
		MaxIdleConns:    config.MaxIdleConns,
	}
	return &http.Client{
		Timeout:   config.HTTPTimeout,
		Transport: tr,
	}
}
