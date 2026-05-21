package skins

import (
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type collectorReader struct {
	price      *NewSkin
	err        error
	currencies []string
	mu         sync.Mutex
}

func (r *collectorReader) GetByMarketHashName(_ context.Context, marketHashName, currency string) (*NewSkin, error) {
	r.mu.Lock()
	r.currencies = append(r.currencies, currency)
	r.mu.Unlock()

	if r.err != nil {
		return nil, r.err
	}
	price := *r.price
	price.MarketHashName = marketHashName
	return &price, nil
}

func (r *collectorReader) lastCurrency() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.currencies) == 0 {
		return ""
	}
	return r.currencies[len(r.currencies)-1]
}

type collectorPageBuilder string

func (b collectorPageBuilder) BuildMarketPageURL(marketHashName string) string {
	return string(b) + marketHashName
}

func TestDefaultSavedSkinPriceCollectorKeepsPartialSuccessAndFallbackPages(t *testing.T) {
	now := time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC)
	steam := &collectorReader{price: &NewSkin{PriceText: "$12.50", PageURL: "steam-page"}}
	lis := &collectorReader{price: &NewSkin{PriceText: "$11.90"}}
	cstm := &collectorReader{err: errors.New("bad status 502")}

	result, err := DefaultSavedSkinPriceCollector{
		Sources: []PriceSource{
			ReaderPriceSource{SourceID: "steam", SourceLabel: "Steam", Reader: steam, Now: func() time.Time { return now }},
			ReaderPriceSource{SourceID: "lisskins", SourceLabel: "LisSkins", Reader: lis, Pages: collectorPageBuilder("lis/"), CurrencyOverride: LisSkinsCurrency, Now: func() time.Time { return now }},
			ReaderPriceSource{SourceID: "cstm", SourceLabel: "CS TM", Reader: cstm, Pages: collectorPageBuilder("cstm/"), Now: func() time.Time { return now }},
		},
		Now: func() time.Time { return now },
	}.Collect(context.Background(), SavedSkin{
		MarketHashName: "AK-47 | Redline",
		CSTMPageURL:    "existing-cstm",
		CSTMPriceText:  "old-cstm",
	}, UpdateSavedSkinPriceParams{MarketHashName: "AK-47 | Redline", Currency: "3"})
	if err != nil {
		t.Fatalf("collect prices: %v", err)
	}

	if result.SteamPriceText != "$12.50" || result.SteamPageURL != "steam-page" {
		t.Fatalf("unexpected steam result: %+v", result)
	}
	if result.LisSkinsPriceText != "$11.90" || result.LisSkinsPageURL != "lis/AK-47 | Redline" {
		t.Fatalf("unexpected lisskins result: %+v", result)
	}
	if result.CSTMPriceText != "old-cstm" || result.CSTMPageURL != "existing-cstm" {
		t.Fatalf("expected failed cstm source to preserve saved values, got %+v", result)
	}
	if result.SteamUpdatedAt != now || result.LisSkinsUpdatedAt != now {
		t.Fatalf("expected successful sources to use injected clock, got %+v", result)
	}
	if result.Currency != "3" {
		t.Fatalf("expected canonical selected currency, got %q", result.Currency)
	}
	if len(result.Prices) != 2 {
		t.Fatalf("expected two successful snapshots, got %+v", result.Prices)
	}
	if steam.lastCurrency() != "3" {
		t.Fatalf("expected steam to use selected currency, got %q", steam.lastCurrency())
	}
	if lis.lastCurrency() != LisSkinsCurrency {
		t.Fatalf("expected lisskins to always use USD, got %q", lis.lastCurrency())
	}
}

func TestDefaultSavedSkinPriceCollectorReturnsSourceErrorWhenAllSourcesFail(t *testing.T) {
	_, err := DefaultSavedSkinPriceCollector{
		Sources: []PriceSource{
			ReaderPriceSource{SourceID: "steam", SourceLabel: "Steam", Reader: &collectorReader{err: errors.New("context deadline exceeded")}},
		},
	}.Collect(context.Background(), SavedSkin{}, UpdateSavedSkinPriceParams{MarketHashName: "AK-47 | Redline", Currency: "1"})
	if err == nil {
		t.Fatalf("expected collect error")
	}
	if !errx.IsCode(err, errx.CodeTimeout) {
		t.Fatalf("expected timeout error, got %v", err)
	}
}
