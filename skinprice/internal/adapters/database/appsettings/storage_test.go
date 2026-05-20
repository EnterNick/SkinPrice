package appsettings

import (
	"SkinPrice/skinprice/internal/adapters/database"
	settings "SkinPrice/skinprice/internal/application/settings"
	"testing"
)

func TestStoragePersistsSavedSkinsViewMode(t *testing.T) {
	connection, err := database.New(&database.Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })

	if err := database.EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	storage := Storage{Conn: connection}
	expected := settings.AppSettings{
		Currency:                   "3",
		AutoRefreshIntervalSeconds: 45,
		SavedSkinsViewMode:         "cards",
	}

	if err := storage.SaveAppSettings(expected); err != nil {
		t.Fatalf("save app settings: %v", err)
	}

	actual, err := storage.GetAppSettings()
	if err != nil {
		t.Fatalf("get app settings: %v", err)
	}

	if actual.Currency != expected.Currency {
		t.Fatalf("expected currency %q, got %q", expected.Currency, actual.Currency)
	}
	if actual.AutoRefreshIntervalSeconds != expected.AutoRefreshIntervalSeconds {
		t.Fatalf("expected auto refresh %d, got %d", expected.AutoRefreshIntervalSeconds, actual.AutoRefreshIntervalSeconds)
	}
	if actual.SavedSkinsViewMode != expected.SavedSkinsViewMode {
		t.Fatalf("expected view mode %q, got %q", expected.SavedSkinsViewMode, actual.SavedSkinsViewMode)
	}
}
