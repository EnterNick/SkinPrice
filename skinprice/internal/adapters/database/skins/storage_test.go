package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"context"
	"sync"
	"testing"
)

type fakeMarketStorage struct {
	prices      map[string]string
	priceCents  map[string]int64
	errors      map[string]error
	currencies  map[string]string
	pageURLs    map[string]string
	beforeFetch func()
}

func (f fakeMarketStorage) GetByMarketHashName(marketHashName, currency string) (*appskins.NewSkin, error) {
	if f.beforeFetch != nil {
		f.beforeFetch()
	}
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
		PageURL:        f.pageURLs[marketHashName],
	}, nil
}

func priceCentsPtr(value int64, ok bool) *int64 {
	if !ok {
		return nil
	}
	return &value
}

func newTestStorage(t *testing.T, steamReader marketPriceReader, lisSkinsReader marketPriceReader) *Storage {
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
		Conn:            connection,
		SteamStorage:    steamReader,
		LisSkinsStorage: lisSkinsReader,
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
	storage := newTestStorage(t, fakeMarketStorage{}, fakeMarketStorage{})

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
	storage := newTestStorage(t, fakeMarketStorage{}, fakeMarketStorage{})
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
	if list.Items[0].SteamPageURL != "page" {
		t.Fatalf("expected steam page url fallback from saved page, got %q", list.Items[0].SteamPageURL)
	}
	if list.Items[0].LisSkinsPageURL == "" {
		t.Fatalf("expected lisskins page url to be initialized")
	}
}

func TestUpdateSavedSkinPriceUpdatesStoredValue(t *testing.T) {
	currencies := map[string]string{}
	storage := newTestStorage(t, fakeMarketStorage{
		prices:     map[string]string{"AK-47 | Redline": "$12.50"},
		priceCents: map[string]int64{"AK-47 | Redline": 1250},
		currencies: currencies,
		pageURLs:   map[string]string{"AK-47 | Redline": "steam-page"},
	}, fakeMarketStorage{
		prices:   map[string]string{"AK-47 | Redline": "$11.90"},
		pageURLs: map[string]string{"AK-47 | Redline": "lis-page"},
	})
	saveFixtureSkin(t, storage, "AK-47 | Redline")

	result, err := storage.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
		MarketHashName: "AK-47 | Redline",
		Currency:       "1",
	})
	if err != nil {
		t.Fatalf("update price: %v", err)
	}
	if result.SteamPriceText != "$12.50" {
		t.Fatalf("expected updated steam price, got %q", result.SteamPriceText)
	}
	if result.LisSkinsPriceText != "$11.90" {
		t.Fatalf("expected updated lis price, got %q", result.LisSkinsPriceText)
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
	if list.Items[0].SteamPriceText != "$12.50" {
		t.Fatalf("expected saved steam price text to be updated, got %q", list.Items[0].SteamPriceText)
	}
	if list.Items[0].LisSkinsPriceText != "$11.90" {
		t.Fatalf("expected saved lis price text to be updated, got %q", list.Items[0].LisSkinsPriceText)
	}
	if list.Items[0].SteamUpdatedAt.IsZero() {
		t.Fatalf("expected steam updated_at to be set")
	}
	if list.Items[0].LisSkinsUpdatedAt.IsZero() {
		t.Fatalf("expected lisskins updated_at to be set")
	}
	if list.Items[0].Currency != "1" {
		t.Fatalf("expected saved canonical currency, got %q", list.Items[0].Currency)
	}
}

func TestUpdateSavedSkinPriceFormatsPriceBySelectedCurrency(t *testing.T) {
	lisCurrencies := map[string]string{}
	storage := newTestStorage(t, fakeMarketStorage{
		prices: map[string]string{
			"AK-47 | Redline": "$12.50",
			"M4A4 | Asiimov":  "$9.99",
		},
		priceCents: map[string]int64{
			"AK-47 | Redline": 1250,
			"M4A4 | Asiimov":  999,
		},
	}, fakeMarketStorage{
		prices:     map[string]string{"M4A4 | Asiimov": "$8.88"},
		priceCents: map[string]int64{"M4A4 | Asiimov": 888},
		currencies: lisCurrencies,
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
	if usdResult.SteamPriceText != "$12.50" {
		t.Fatalf("expected formatted USD price, got %q", usdResult.SteamPriceText)
	}

	eurResult, err := storage.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
		MarketHashName: "M4A4 | Asiimov",
		Currency:       "3",
	})
	if err != nil {
		t.Fatalf("update EUR price: %v", err)
	}
	if eurResult.SteamPriceText != "€9.99" {
		t.Fatalf("expected formatted EUR price, got %q", eurResult.SteamPriceText)
	}
	if eurResult.LisSkinsPriceText != "$8.88" {
		t.Fatalf("expected lisskins price to stay in USD, got %q", eurResult.LisSkinsPriceText)
	}
	if lisCurrencies["M4A4 | Asiimov"] != "1" {
		t.Fatalf("expected lisskins updater to always receive USD currency, got %q", lisCurrencies["M4A4 | Asiimov"])
	}
}

func TestUpdateAllSavedSkinsPricesReturnsPartialSuccess(t *testing.T) {
	currencies := map[string]string{}
	storage := newTestStorage(t, fakeMarketStorage{
		prices:     map[string]string{"AK-47 | Redline": "$12.50"},
		errors:     map[string]error{"M4A4 | Asiimov": appskins.ErrNewSkinsRequestFailed},
		currencies: currencies,
	}, fakeMarketStorage{
		prices: map[string]string{"AK-47 | Redline": "$11.90"},
		errors: map[string]error{"M4A4 | Asiimov": appskins.ErrNewSkinsRequestFailed},
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

	list, err := storage.GetSavedList(&application.Pagination{Limit: 20, Offset: 0})
	if err != nil {
		t.Fatalf("get saved list after bulk update: %v", err)
	}
	if list.Items[1].LisSkinsPriceText != "$11.90" {
		t.Fatalf("expected lisskins price to update during bulk refresh, got %q", list.Items[1].LisSkinsPriceText)
	}
}

func TestUpdateSavedSkinPricePreservesSuccessfulSourceOnPartialFailure(t *testing.T) {
	storage := newTestStorage(t, fakeMarketStorage{
		prices:   map[string]string{"AK-47 | Redline": "$12.50"},
		pageURLs: map[string]string{"AK-47 | Redline": "steam-page"},
	}, fakeMarketStorage{
		errors: map[string]error{"AK-47 | Redline": appskins.ErrNewSkinsRequestFailed},
	})
	saveFixtureSkin(t, storage, "AK-47 | Redline")

	result, err := storage.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
		MarketHashName: "AK-47 | Redline",
		Currency:       "1",
	})
	if err != nil {
		t.Fatalf("expected partial update to succeed, got %v", err)
	}
	if result.SteamPriceText != "$12.50" {
		t.Fatalf("expected steam update to be preserved, got %q", result.SteamPriceText)
	}
	if result.LisSkinsPriceText != "" {
		t.Fatalf("expected failed lisskins update to keep empty price, got %q", result.LisSkinsPriceText)
	}
}

func TestUpdateSavedSkinPriceFetchesSourcesConcurrently(t *testing.T) {
	var mu sync.Mutex
	started := 0
	release := make(chan struct{})
	bothStarted := make(chan struct{})

	markStarted := func() {
		mu.Lock()
		defer mu.Unlock()
		started++
		if started == 2 {
			close(bothStarted)
		}
	}

	blockingStorage := fakeMarketStorage{
		prices: map[string]string{"AK-47 | Redline": "$12.50"},
		beforeFetch: func() {
			markStarted()
			<-release
		},
	}

	storage := newTestStorage(t, blockingStorage, blockingStorage)
	saveFixtureSkin(t, storage, "AK-47 | Redline")

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = storage.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
			MarketHashName: "AK-47 | Redline",
			Currency:       "1",
		})
	}()

	<-bothStarted
	close(release)
	<-done
}

func TestGetSavedListNormalizesLegacyCurrencyValues(t *testing.T) {
	storage := newTestStorage(t, fakeMarketStorage{}, fakeMarketStorage{})
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
	storage := newTestStorage(t, fakeMarketStorage{}, fakeMarketStorage{})
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
