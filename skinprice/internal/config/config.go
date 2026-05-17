package config

import (
	"SkinPrice/skinprice/internal/shared/utils"
	"time"
)

type Config struct {
	AppEnv          string
	SteamBaseURL    string
	LisSkinsBaseURL string
	HTTPTimeout     time.Duration
	CacheTTL        time.Duration
	MaxIdleConns    int
}

func Load() Config {
	return Config{
		AppEnv:          utils.GetStrWDefault("APP_ENV", "local"),
		SteamBaseURL:    utils.GetStrWDefault("STEAM_BASE_URL", "https://steamcommunity.com/market"),
		LisSkinsBaseURL: utils.GetStrWDefault("LISSKINS_BASE_URL", "https://api.lis-skins.ru/v1"),
		HTTPTimeout:     time.Duration(utils.GetIntWDefault("HTTP_TIMEOUT_SECONDS", 10)) * time.Second,
		CacheTTL:        time.Duration(utils.GetIntWDefault("CACHE_TTL_SECONDS", 300)) * time.Second,
		MaxIdleConns:    utils.GetIntWDefault("MAX_IDLE_CONNS", 10),
	}
}
