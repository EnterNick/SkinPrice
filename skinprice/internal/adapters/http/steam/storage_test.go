package steam

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetByMarketHashNameUsesPriceOverviewText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/priceoverview/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("appid") != "730" {
			t.Fatalf("unexpected appid: %s", r.URL.Query().Get("appid"))
		}
		if r.URL.Query().Get("market_hash_name") != "Gamma Case Key" {
			t.Fatalf("unexpected market_hash_name: %s", r.URL.Query().Get("market_hash_name"))
		}
		if r.URL.Query().Get("currency") != "5" {
			t.Fatalf("unexpected currency: %s", r.URL.Query().Get("currency"))
		}

		_, _ = w.Write([]byte(`{"success":true,"lowest_price":"1496,28 руб."}`))
	}))
	defer server.Close()

	storage := &Storage{
		Client:  server.Client(),
		BaseURL: server.URL,
	}

	result, err := storage.GetByMarketHashName("Gamma Case Key", "5")
	if err != nil {
		t.Fatalf("GetByMarketHashName() error = %v", err)
	}
	if result.PriceText != "1496,28 руб." {
		t.Fatalf("expected RUB price text, got %q", result.PriceText)
	}
	if result.PageURL != server.URL+"/listings/730/"+url.PathEscape("Gamma Case Key") {
		t.Fatalf("unexpected page url: %q", result.PageURL)
	}
}

func TestGetByMarketHashNameFallsBackToMedianPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"median_price":"17,56€"}`))
	}))
	defer server.Close()

	storage := &Storage{
		Client:  server.Client(),
		BaseURL: server.URL,
	}

	result, err := storage.GetByMarketHashName("Gamma Case Key", "3")
	if err != nil {
		t.Fatalf("GetByMarketHashName() error = %v", err)
	}
	if result.PriceText != "17,56€" {
		t.Fatalf("expected EUR median price text, got %q", result.PriceText)
	}
}
