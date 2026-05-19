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
	LogLevel        string
	LogFormat       string
	LogToFile       bool
	LogFilePath     string
	LogMaxSizeMB    int
	LogMaxBackups   int
	LogMaxAgeDays   int
	LogCompress     bool
}

func Load() Config {
	return Config{
		AppEnv:          utils.GetStrWDefault("APP_ENV", "local"),
		SteamBaseURL:    utils.GetStrWDefault("STEAM_BASE_URL", "https://steamcommunity.com/market"),
		LisSkinsBaseURL: utils.GetStrWDefault("LISSKINS_BASE_URL", "https://api.lis-skins.ru/v1"),
		HTTPTimeout:     time.Duration(utils.GetIntWDefault("HTTP_TIMEOUT_SECONDS", 10)) * time.Second,
		CacheTTL:        time.Duration(utils.GetIntWDefault("CACHE_TTL_SECONDS", 300)) * time.Second,
		MaxIdleConns:    utils.GetIntWDefault("MAX_IDLE_CONNS", 10),
		LogLevel:        utils.GetStrWDefault("LOG_LEVEL", defaultLogLevel()),
		LogFormat:       utils.GetStrWDefault("LOG_FORMAT", "text"),
		LogToFile:       utils.GetBoolWDefault("LOG_TO_FILE", true),
		LogFilePath:     utils.GetStrWDefault("LOG_FILE_PATH", ""),
		LogMaxSizeMB:    utils.GetIntWDefault("LOG_MAX_SIZE_MB", 20),
		LogMaxBackups:   utils.GetIntWDefault("LOG_MAX_BACKUPS", 5),
		LogMaxAgeDays:   utils.GetIntWDefault("LOG_MAX_AGE_DAYS", 14),
		LogCompress:     utils.GetBoolWDefault("LOG_COMPRESS", true),
	}
}

func defaultLogLevel() string {
	if utils.GetStrWDefault("APP_ENV", "local") == "local" {
		return "debug"
	}
	return "info"
}
