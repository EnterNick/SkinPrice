package steam

import (
	"SkinPrice/skinprice/internal/application"
	"SkinPrice/skinprice/internal/application/skins"
	"context"
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

	result, err := storage.GetList(context.Background(), skins.SearchCriteria{}, &application.Pagination{Limit: 1, Offset: 0})
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

func TestBuildSteamMarketSearchParamsIncludesServerFilters(t *testing.T) {
	priceMin := "100"
	priceMax := "2500"

	params := buildSteamMarketSearchParams(
		skins.SearchCriteria{
			SortColumn:         "price",
			SortDir:            "asc",
			PriceMin:           &priceMin,
			PriceMax:           &priceMax,
			SearchDescriptions: true,
			Types:              []string{"tag_CSGO_Type_Rifle"},
			Weapons:            []string{"tag_weapon_ak47", "tag_weapon_awp"},
			Exteriors:          []string{"tag_WearCategory0"},
		},
		&application.Pagination{Limit: 20, Offset: 40},
	)

	if got := params.Get("start"); got != "40" {
		t.Fatalf("expected start=40, got %q", got)
	}
	if got := params.Get("count"); got != "20" {
		t.Fatalf("expected count=20, got %q", got)
	}
	if got := params.Get("sort_column"); got != "price" {
		t.Fatalf("expected sort_column=price, got %q", got)
	}
	if got := params.Get("sort_dir"); got != "asc" {
		t.Fatalf("expected sort_dir=asc, got %q", got)
	}
	if got := params.Get("price_min"); got != "100" {
		t.Fatalf("expected price_min=100, got %q", got)
	}
	if got := params.Get("price_max"); got != "2500" {
		t.Fatalf("expected price_max=2500, got %q", got)
	}
	if got := params.Get("search_descriptions"); got != "1" {
		t.Fatalf("expected search_descriptions=1, got %q", got)
	}
	if got := params["category_730_Type[]"]; len(got) != 1 || got[0] != "tag_CSGO_Type_Rifle" {
		t.Fatalf("unexpected type filters: %#v", got)
	}
	if got := params["category_730_Weapon[]"]; len(got) != 2 || got[0] != "tag_weapon_ak47" || got[1] != "tag_weapon_awp" {
		t.Fatalf("unexpected weapon filters: %#v", got)
	}
	if got := params["category_730_Exterior[]"]; len(got) != 1 || got[0] != "tag_WearCategory0" {
		t.Fatalf("unexpected exterior filters: %#v", got)
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

	result, err := storage.GetByMarketHashName(context.Background(), "Gamma Case Key", "5")
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

	result, err := storage.GetByMarketHashName(context.Background(), "Gamma Case Key", "3")
	if err != nil {
		t.Fatalf("GetByMarketHashName() error = %v", err)
	}
	if result.PriceText != "17,56€" {
		t.Fatalf("expected EUR median price text, got %q", result.PriceText)
	}
}
