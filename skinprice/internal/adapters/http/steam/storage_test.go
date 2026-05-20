package steam

import (
	"SkinPrice/skinprice/internal/application"
	"SkinPrice/skinprice/internal/application/skins"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetListParsesNameColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/render/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{
			"success": true,
			"start": 0,
			"pagesize": 1,
			"total_count": 1,
			"results": [
				{
					"name": "Gamma Doppler",
					"hash_name": "Gamma Doppler",
					"sell_listings": 12,
					"sell_price": 1234,
					"sell_price_text": "$12.34",
					"asset_description": {
						"icon_url": "icon-hash",
						"name_color": "8847ff"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	storage := &Storage{
		Client:  server.Client(),
		BaseURL: server.URL,
	}

	result, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 1, Offset: 0})
	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].NameColor != "8847ff" {
		t.Fatalf("expected name color to be parsed, got %q", result.Items[0].NameColor)
	}
}

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
