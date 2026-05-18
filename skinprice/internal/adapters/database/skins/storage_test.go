package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"context"
	"testing"
)

type fakeSteamStorage struct {
	prices     map[string]string
	priceCents map[string]int64
	errors     map[string]error
	currencies map[string]string
}

func (f fakeSteamStorage) GetByMarketHashName(marketHashName, currency string) (*appskins.NewSkin, error) {
	if err := f.errors[marketHashName]; err != nil {
		return nil, err
	}
	if f.currencies != nil {
		f.currencies[marketHashName] = currency
	}
	return &appskins.NewSkin{
		MarketHashName: marketHashName,
		PriceText:      f.prices[marketHashName],
		PriceCents:     priceCentsPtr(f.priceCents[marketHashName], f.priceCents != nil),
	}, nil
}

func priceCentsPtr(value int64, ok bool) *int64 {
	if !ok {
		return nil
	}
	return &value
}

func newTestStorage(t *testing.T, reader steamPriceReader) *Storage {
	t.Helper()

	connection, err := database.New(&database.Config{
		Driver: "sqlite3",
		DBName: t.TempDir() + "/skinprice.db",
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() {
		_ = connection.Close()
	})

	if err := database.EnsureSchema(connection); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	return &Storage{
		Conn:         connection,
		SteamStorage: reader,
	}
}

func saveFixtureSkin(t *testing.T, storage *Storage, hash string) {
	t.Helper()

	if _, err := storage.Save(appskins.SaveSkinParams{
		MarketHashName: hash,
		DisplayName:    hash,
		IconURL:        "icon",
		PageURL:        "page",
	}); err != nil {
		t.Fatalf("save fixture skin: %v", err)
	}
}

func TestSaveSkinIgnoresDuplicates(t *testing.T) {
	storage := newTestStorage(t, fakeSteamStorage{})

	first, err := storage.Save(appskins.SaveSkinParams{
		MarketHashName: "AK-47 | Redline",
		DisplayName:    "AK-47 | Redline",
		IconURL:        "icon",
		PageURL:        "page",
	})
	if err != nil {
		t.Fatalf("save first skin: %v", err)
	}
	if !first.Created {
		t.Fatalf("expected first save to create a record")
	}

	second, err := storage.Save(appskins.SaveSkinParams{
		MarketHashName: "AK-47 | Redline",
		DisplayName:    "AK-47 | Redline",
		IconURL:        "icon",
		PageURL:        "page",
	})
	if err != nil {
		t.Fatalf("save duplicate skin: %v", err)
	}
	if second.Created {
		t.Fatalf("expected duplicate save to be ignored")
	}
}

func TestGetSavedListReturnsSavedItems(t *testing.T) {
	storage := newTestStorage(t, fakeSteamStorage{})
	saveFixtureSkin(t, storage, "AK-47 | Redline")

	list, err := storage.GetSavedList(&application.Pagination{Limit: 20, Offset: 0})
	if err != nil {
		t.Fatalf("get saved list: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list.Items))
	}
	if list.Items[0].MarketHashName != "AK-47 | Redline" {
		t.Fatalf("unexpected item: %+v", list.Items[0])
	}
	if list.Items[0].Currency != "1" {
		t.Fatalf("expected default canonical currency, got %q", list.Items[0].Currency)
	}
}

func TestUpdateSavedSkinPriceUpdatesStoredValue(t *testing.T) {
	currencies := map[string]string{}
	storage := newTestStorage(t, fakeSteamStorage{
		prices:     map[string]string{"AK-47 | Redline": "$12.50"},
		priceCents: map[string]int64{"AK-47 | Redline": 1250},
		currencies: currencies,
	})
	saveFixtureSkin(t, storage, "AK-47 | Redline")

	result, err := storage.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
		MarketHashName: "AK-47 | Redline",
		Currency:       "1",
	})
	if err != nil {
		t.Fatalf("update price: %v", err)
	}
	if result.PriceText != "$12.50" {
		t.Fatalf("expected updated price, got %q", result.PriceText)
	}
	if result.Currency != "1" {
		t.Fatalf("expected canonical currency code, got %q", result.Currency)
	}
	if currencies["AK-47 | Redline"] != "1" {
		t.Fatalf("expected updater to receive canonical currency, got %q", currencies["AK-47 | Redline"])
	}

	list, err := storage.GetSavedList(&application.Pagination{Limit: 20, Offset: 0})
	if err != nil {
		t.Fatalf("get saved list after update: %v", err)
	}
	if list.Items[0].PriceText != "$12.50" {
		t.Fatalf("expected saved price text to be updated, got %q", list.Items[0].PriceText)
	}
	if list.Items[0].UpdatedAt.IsZero() {
		t.Fatalf("expected updated_at to be set")
	}
	if list.Items[0].Currency != "1" {
		t.Fatalf("expected saved canonical currency, got %q", list.Items[0].Currency)
	}
}

func TestUpdateSavedSkinPriceFormatsPriceBySelectedCurrency(t *testing.T) {
	storage := newTestStorage(t, fakeSteamStorage{
		prices: map[string]string{
			"AK-47 | Redline": "$12.50",
			"M4A4 | Asiimov":  "$9.99",
		},
		priceCents: map[string]int64{
			"AK-47 | Redline": 1250,
			"M4A4 | Asiimov":  999,
		},
	})
	saveFixtureSkin(t, storage, "AK-47 | Redline")
	saveFixtureSkin(t, storage, "M4A4 | Asiimov")

	usdResult, err := storage.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
		MarketHashName: "AK-47 | Redline",
		Currency:       "1",
	})
	if err != nil {
		t.Fatalf("update USD price: %v", err)
	}
	if usdResult.PriceText != "$12.50" {
		t.Fatalf("expected formatted USD price, got %q", usdResult.PriceText)
	}

	eurResult, err := storage.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
		MarketHashName: "M4A4 | Asiimov",
		Currency:       "3",
	})
	if err != nil {
		t.Fatalf("update EUR price: %v", err)
	}
	if eurResult.PriceText != "€9.99" {
		t.Fatalf("expected formatted EUR price, got %q", eurResult.PriceText)
	}
}

func TestUpdateAllSavedSkinsPricesReturnsPartialSuccess(t *testing.T) {
	currencies := map[string]string{}
	storage := newTestStorage(t, fakeSteamStorage{
		prices:     map[string]string{"AK-47 | Redline": "$12.50"},
		errors:     map[string]error{"M4A4 | Asiimov": appskins.ErrNewSkinsRequestFailed},
		currencies: currencies,
	})
	saveFixtureSkin(t, storage, "AK-47 | Redline")
	saveFixtureSkin(t, storage, "M4A4 | Asiimov")

	result, err := storage.UpdateAllSavedSkinsPrices(appskins.UpdateAllSavedSkinsPricesParams{Currency: "1"})
	if err != nil {
		t.Fatalf("update all prices: %v", err)
	}
	if result.UpdatedCount != 1 || result.FailedCount != 1 {
		t.Fatalf("unexpected aggregate result: %+v", result)
	}
	if len(result.Failures) != 1 || result.Failures[0].MarketHashName != "M4A4 | Asiimov" {
		t.Fatalf("unexpected failures: %+v", result.Failures)
	}
	if currencies["AK-47 | Redline"] != "1" {
		t.Fatalf("expected successful update to use canonical currency, got %+v", currencies)
	}
}

func TestGetSavedListNormalizesLegacyCurrencyValues(t *testing.T) {
	storage := newTestStorage(t, fakeSteamStorage{})
	saveFixtureSkin(t, storage, "AK-47 | Redline")

	if _, err := storage.Conn.DB().ExecContext(context.Background(), `UPDATE skins SET currency = 'USD' WHERE market_hash_name = ?`, "AK-47 | Redline"); err != nil {
		t.Fatalf("set legacy currency: %v", err)
	}

	list, err := storage.GetSavedList(&application.Pagination{Limit: 20, Offset: 0})
	if err != nil {
		t.Fatalf("get saved list: %v", err)
	}
	if list.Items[0].Currency != "1" {
		t.Fatalf("expected legacy currency to normalize to canonical code, got %q", list.Items[0].Currency)
	}
}

func TestDeleteSavedSkinRemovesRecord(t *testing.T) {
	storage := newTestStorage(t, fakeSteamStorage{})
	saveFixtureSkin(t, storage, "AK-47 | Redline")

	if err := storage.DeleteSavedSkin(appskins.DeleteSavedSkinParams{MarketHashName: "AK-47 | Redline"}); err != nil {
		t.Fatalf("delete saved skin: %v", err)
	}

	list, err := storage.GetSavedList(&application.Pagination{Limit: 20, Offset: 0})
	if err != nil {
		t.Fatalf("get saved list after delete: %v", err)
	}
	if len(list.Items) != 0 {
		t.Fatalf("expected empty list after delete, got %d items", len(list.Items))
	}
}
