package lisskins

import (
	"SkinPrice/skinprice/internal/config"
	"net/http"
)

func NewLisSkinsClient(config config.Config) *http.Client {
	tr := &http.Transport{
		IdleConnTimeout: config.HTTPTimeout,
		MaxIdleConns:    config.MaxIdleConns,
	}
	return &http.Client{
		Transport: tr,
	}
}
