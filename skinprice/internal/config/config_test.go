package config

import (
	"encoding/base64"
	"testing"
)

func testKey() string {
	return base64.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012"))
}

func TestLoadReadsLoggingConfig(t *testing.T) {
	t.Setenv("APP_ENV", "prod")
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("LOG_FORMAT", "json")
	t.Setenv("LOG_TO_FILE", "false")
	t.Setenv("LOG_FILE_PATH", "/tmp/skinprice.log")
	t.Setenv("LOG_MAX_SIZE_MB", "42")
	t.Setenv("LOG_MAX_BACKUPS", "7")
	t.Setenv("LOG_MAX_AGE_DAYS", "30")
	t.Setenv("LOG_COMPRESS", "false")
	t.Setenv("TOKEN_ENCRYPTION_KEY", testKey())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.LogLevel != "warn" {
		t.Fatalf("LogLevel = %q, want warn", cfg.LogLevel)
	}
}

func TestLoadDefaultsDebugLevelInLocalEnv(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("TOKEN_ENCRYPTION_KEY", testKey())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want debug", cfg.LogLevel)
	}
}

func TestLoadFailsWithoutTokenKey(t *testing.T) {
	t.Setenv("TOKEN_ENCRYPTION_KEY", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() expected error when TOKEN_ENCRYPTION_KEY is missing")
	}
}
