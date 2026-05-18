package config

import "testing"

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

	cfg := Load()

	if cfg.LogLevel != "warn" {
		t.Fatalf("LogLevel = %q, want warn", cfg.LogLevel)
	}
	if cfg.LogFormat != "json" {
		t.Fatalf("LogFormat = %q, want json", cfg.LogFormat)
	}
	if cfg.LogToFile {
		t.Fatal("LogToFile = true, want false")
	}
	if cfg.LogFilePath != "/tmp/skinprice.log" {
		t.Fatalf("LogFilePath = %q, want /tmp/skinprice.log", cfg.LogFilePath)
	}
	if cfg.LogMaxSizeMB != 42 {
		t.Fatalf("LogMaxSizeMB = %d, want 42", cfg.LogMaxSizeMB)
	}
	if cfg.LogMaxBackups != 7 {
		t.Fatalf("LogMaxBackups = %d, want 7", cfg.LogMaxBackups)
	}
	if cfg.LogMaxAgeDays != 30 {
		t.Fatalf("LogMaxAgeDays = %d, want 30", cfg.LogMaxAgeDays)
	}
	if cfg.LogCompress {
		t.Fatal("LogCompress = true, want false")
	}
}

func TestLoadDefaultsDebugLevelInLocalEnv(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("LOG_LEVEL", "")

	cfg := Load()

	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want debug", cfg.LogLevel)
	}
}
