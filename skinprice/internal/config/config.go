package config

import (
	"SkinPrice/skinprice/internal/shared/utils"
	"encoding/base64"
	"fmt"
	"time"
)

type Config struct {
	AppEnv               string
	SteamBaseURL         string
	LisSkinsBaseURL      string
	HTTPTimeout          time.Duration
	BulkPriceUpdateDelay time.Duration
	CacheTTL             time.Duration
	MaxIdleConns         int
	LogLevel             string
	LogFormat            string
	LogToFile            bool
	LogFilePath          string
	LogMaxSizeMB         int
	LogMaxBackups        int
	LogMaxAgeDays        int
	LogCompress          bool
	TokenEncryptionKey   string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:               utils.GetStrWDefault("APP_ENV", "local"),
		SteamBaseURL:         utils.GetStrWDefault("STEAM_BASE_URL", "https://steamcommunity.com/market"),
		LisSkinsBaseURL:      utils.GetStrWDefault("LISSKINS_BASE_URL", "https://api.lis-skins.com/v1"),
		HTTPTimeout:          time.Duration(utils.GetIntWDefault("HTTP_TIMEOUT_SECONDS", 10)) * time.Second,
		BulkPriceUpdateDelay: time.Duration(utils.GetIntWDefault("BULK_PRICE_UPDATE_DELAY_MS", 1200)) * time.Millisecond,
		CacheTTL:             time.Duration(utils.GetIntWDefault("CACHE_TTL_SECONDS", 300)) * time.Second,
		MaxIdleConns:         utils.GetIntWDefault("MAX_IDLE_CONNS", 10),
		LogLevel:             utils.GetStrWDefault("LOG_LEVEL", defaultLogLevel()),
		LogFormat:            utils.GetStrWDefault("LOG_FORMAT", "text"),
		LogToFile:            utils.GetBoolWDefault("LOG_TO_FILE", true),
		LogFilePath:          utils.GetStrWDefault("LOG_FILE_PATH", "../../logs/skinprice.log"),
		LogMaxSizeMB:         utils.GetIntWDefault("LOG_MAX_SIZE_MB", 20),
		LogMaxBackups:        utils.GetIntWDefault("LOG_MAX_BACKUPS", 5),
		LogMaxAgeDays:        utils.GetIntWDefault("LOG_MAX_AGE_DAYS", 14),
		LogCompress:          utils.GetBoolWDefault("LOG_COMPRESS", true),
		TokenEncryptionKey:   utils.GetStrWDefault("TOKEN_ENCRYPTION_KEY", ""),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.TokenEncryptionKey == "" {
		return nil
	}
	decoded, err := base64.StdEncoding.DecodeString(c.TokenEncryptionKey)
	if err != nil {
		return fmt.Errorf("TOKEN_ENCRYPTION_KEY must be valid base64: %w", err)
	}
	if len(decoded) != 32 {
		return fmt.Errorf("TOKEN_ENCRYPTION_KEY must decode to 32 bytes (got %d)", len(decoded))
	}
	return nil
}

func defaultLogLevel() string {
	if utils.GetStrWDefault("APP_ENV", "local") == "local" {
		return "debug"
	}
	return "info"
}
