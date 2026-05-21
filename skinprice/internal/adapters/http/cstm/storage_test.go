package cstm

import (
	"SkinPrice/skinprice/internal/application/skins"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetByMarketHashNameParsesPriceListAndRespectsCurrency(t *testing.T) {
	var requestedPath atomic.Value
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath.Store(r.URL.Path)
		_, _ = w.Write([]byte(`[{"market_hash_name":"AK-47 | Redline","volume":"10","price":"125.25"}]`))
	}))
	defer server.Close()

	storage := &Storage{
		Client:         server.Client(),
		BaseURL:        server.URL,
		RequestTimeout: time.Second,
		CacheTTL:       time.Minute,
	}

	result, err := storage.GetByMarketHashName("AK-47 | Redline", "5")
	if err != nil {
		t.Fatalf("lookup price: %v", err)
	}
	if got := requestedPath.Load(); got != "/api/v2/prices/RUB.json" {
		t.Fatalf("unexpected path: %v", got)
	}
	if result.PriceCents == nil || *result.PriceCents != 12525 {
		t.Fatalf("unexpected cents: %+v", result.PriceCents)
	}
	if result.PageURL != BuildMarketPageURL(server.URL, "AK-47 | Redline") {
		t.Fatalf("unexpected page url: %q", result.PageURL)
	}
	if result.PriceText != "125.25 ₽" {
		t.Fatalf("unexpected price text: %q", result.PriceText)
	}
}

func TestGetByMarketHashNameUsesCacheWithinTTL(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_, _ = w.Write([]byte(`[{"market_hash_name":"AK-47 | Redline","volume":"10","price":"125.25"}]`))
	}))
	defer server.Close()

	storage := &Storage{
		Client:         server.Client(),
		BaseURL:        server.URL,
		RequestTimeout: time.Second,
		CacheTTL:       time.Minute,
	}

	for i := 0; i < 2; i++ {
		if _, err := storage.GetByMarketHashName("AK-47 | Redline", "5"); err != nil {
			t.Fatalf("lookup %d: %v", i, err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("expected single upstream fetch, got %d", calls.Load())
	}
}

func TestGetByMarketHashNameParsesWrappedObjectResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[{"market_hash_name":"AK-47 | Redline","volume":"10","price":"125.25"}]}`))
	}))
	defer server.Close()

	storage := &Storage{
		Client:         server.Client(),
		BaseURL:        server.URL,
		RequestTimeout: time.Second,
		CacheTTL:       time.Minute,
	}

	result, err := storage.GetByMarketHashName("AK-47 | Redline", "5")
	if err != nil {
		t.Fatalf("lookup wrapped response: %v", err)
	}
	if result.PriceCents == nil || *result.PriceCents != 12525 {
		t.Fatalf("unexpected cents: %+v", result.PriceCents)
	}
}

func TestGetByMarketHashNameReturnsNotFoundForMissingSkin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"market_hash_name":"M4A4 | Asiimov","volume":"10","price":"125.25"}]`))
	}))
	defer server.Close()

	storage := &Storage{
		Client:         server.Client(),
		BaseURL:        server.URL,
		RequestTimeout: time.Second,
		CacheTTL:       time.Minute,
	}

	_, err := storage.GetByMarketHashName("AK-47 | Redline", "5")
	if !errors.Is(err, skins.ErrNewSkinsResponseUnsuccess) {
		t.Fatalf("unexpected error: %v", err)
	}
}
