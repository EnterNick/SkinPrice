package lisskins

import (
	"SkinPrice/skinprice/internal/application"
	"SkinPrice/skinprice/internal/application/skins"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type tokenProviderStub struct {
	token string
	err   error
}

func (s tokenProviderStub) Execute() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.token, nil
}

func TestGetListAddsAuthorizationHeader(t *testing.T) {
	token := "secret-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := buildLisSkinsAuthHeaderValue(token)
		if got := r.Header.Get(lisSkinsAuthHeader); got != expected {
			t.Fatalf("unexpected auth header: %q", got)
		}
		_, _ = w.Write([]byte(`{"data":[],"total_count":0}`))
	}))
	defer server.Close()

	storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: token}}
	_, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
}

func TestGetListWithMissingTokenSendsRequestWithoutAuthorization(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		if got := r.Header.Get(lisSkinsAuthHeader); got != "" {
			t.Fatalf("expected no auth header, got %q", got)
		}
		_, _ = w.Write([]byte(`{"data":[],"total_count":0}`))
	}))
	defer server.Close()

	storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: ""}}
	_, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("expected request without token to succeed, got %v", err)
	}
	if atomic.LoadInt32(&requests) != 1 {
		t.Fatalf("expected one outbound request, got %d", requests)
	}
}

func TestGetListMapsAuthStatusToInvalidToken(t *testing.T) {
	tests := []int{http.StatusUnauthorized, http.StatusForbidden}
	for _, status := range tests {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer server.Close()

			storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: "token"}}
			_, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10, Offset: 0})
			if !errors.Is(err, skins.ErrLisSkinsTokenInvalid) {
				t.Fatalf("expected ErrLisSkinsTokenInvalid, got %v", err)
			}
		})
	}
}

func TestGetListWithoutTokenLeavesUnauthorizedAsBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: ""}}
	_, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10, Offset: 0})
	if errors.Is(err, skins.ErrLisSkinsTokenInvalid) {
		t.Fatalf("expected unauthorized without token to avoid invalid-token classification")
	}
	if !errors.Is(err, skins.ErrNewSkinsRequestBadStatus) {
		t.Fatalf("expected bad status error, got %v", err)
	}
}

func TestBuildLisSkinsMarketSearchParamsUsesGameAndCursor(t *testing.T) {
	q := buildLisSkinsMarketSearchParams(
		skins.SearchCriteria{},
		&application.Pagination{Limit: 20, Cursor: "next-page-token"},
		"",
	)

	if got := q.Get("game"); got != "csgo" {
		t.Fatalf("expected game=csgo, got %q", got)
	}
	if got := q.Get("app_id"); got != "" {
		t.Fatalf("expected app_id to be omitted, got %q", got)
	}
	if got := q.Get("cursor"); got != "next-page-token" {
		t.Fatalf("expected cursor to be passed through, got %q", got)
	}
	if got := q.Get("names[]"); got != "" {
		t.Fatalf("expected names[] to be empty when no search term is provided, got %q", got)
	}
}

func TestBuildLisSkinsMarketSearchParamsUsesNamesArray(t *testing.T) {
	name := "AK-47 | Redline"
	q := buildLisSkinsMarketSearchParams(
		skins.SearchCriteria{MarketHashName: &name},
		&application.Pagination{Limit: 20},
		"",
	)

	if got := q.Get("names[]"); got != name {
		t.Fatalf("expected names[]=%q, got %q", name, got)
	}
	if got := q.Get("name"); got != "" {
		t.Fatalf("expected legacy name param to be omitted, got %q", got)
	}
}

func TestGetListParsesMetaNextCursor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("game"); got != "csgo" {
			t.Fatalf("expected game=csgo, got %q", got)
		}
		_, _ = w.Write([]byte(`{"data":[{"name":"StatTrak™ MP9 | Featherweight (Field-Tested)","price":0.12,"count":7,"asset_description":{"icon_url":"icon"}}],"meta":{"next_cursor":"cursor-2","count":1}}`))
	}))
	defer server.Close()

	storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: "token"}}
	list, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10})
	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
	if list.NextCursor != "cursor-2" {
		t.Fatalf("expected next cursor to be parsed, got %q", list.NextCursor)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list.Items))
	}
	if list.Items[0].MarketHashName != "StatTrak™ MP9 | Featherweight (Field-Tested)" {
		t.Fatalf("expected market hash name fallback to item name, got %q", list.Items[0].MarketHashName)
	}
	if list.Items[0].PriceCents == nil || *list.Items[0].PriceCents != 12 {
		t.Fatalf("expected decimal price to convert to 12 cents, got %v", list.Items[0].PriceCents)
	}
	if list.Items[0].PriceText != "$0.12" {
		t.Fatalf("expected fallback price text $0.12, got %q", list.Items[0].PriceText)
	}
	if list.Items[0].PageURL != "https://lis-skins.com/market/csgo/stattrak-mp9-featherweight-field-tested/" {
		t.Fatalf("expected lis-skins slug page url, got %q", list.Items[0].PageURL)
	}
}

func TestBuildMarketPageURL(t *testing.T) {
	got := BuildMarketPageURL("StatTrak™ MP9 | Featherweight (Field-Tested)")
	if got != "https://lis-skins.com/market/csgo/stattrak-mp9-featherweight-field-tested/" {
		t.Fatalf("unexpected page url: %q", got)
	}
}

func TestGetByMarketHashNameMatchesNormalizedNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"name":"StatTrak MP9 | Featherweight (Field-Tested)","price":12.34}]}`))
	}))
	defer server.Close()

	storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: ""}}
	skin, err := storage.GetByMarketHashName("StatTrak™ MP9 | Featherweight (Field-Tested)", "1")
	if err != nil {
		t.Fatalf("GetByMarketHashName() error = %v", err)
	}
	if skin == nil {
		t.Fatalf("expected skin result")
	}
	if skin.PriceText != "$12.34" {
		t.Fatalf("expected normalized match to return price, got %q", skin.PriceText)
	}
}

func TestBuildLisSkinsItemSlug(t *testing.T) {
	got := buildLisSkinsItemSlug("StatTrak™ MP9 | Featherweight (Field-Tested)")
	if got != "stattrak-mp9-featherweight-field-tested" {
		t.Fatalf("unexpected slug: %q", got)
	}
}
