package database

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadConfig_DefaultSQLitePathWhenDBNameMissing(t *testing.T) {
	t.Setenv("APP_DB_DRIVER", "sqlite3")
	t.Setenv("APP_DB_NAME", "")
	t.Setenv("APP_DB_PATH", "")
	t.Setenv("LOCALAPPDATA", "")

	cfg := LoadConfig()

	switch runtime.GOOS {
	case "windows":
		if cfg.DBName != "skinprice.db" {
			t.Fatalf("expected fallback db name for windows without LOCALAPPDATA, got %q", cfg.DBName)
		}
	case "darwin":
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)
		cfg = LoadConfig()
		want := filepath.Join(homeDir, "Library", "Application Support", "SkinPrice", "skinprice.db")
		if cfg.DBName != want {
			t.Fatalf("expected %q, got %q", want, cfg.DBName)
		}
	default:
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)
		cfg = LoadConfig()
		want := filepath.Join(homeDir, ".local", "share", "SkinPrice", "skinprice.db")
		if cfg.DBName != want {
			t.Fatalf("expected %q, got %q", want, cfg.DBName)
		}
	}
}
