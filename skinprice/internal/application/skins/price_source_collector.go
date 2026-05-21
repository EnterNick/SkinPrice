package skins

import (
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"strings"
	"sync"
	"time"
)

type DefaultSavedSkinPriceCollector struct {
	Sources     []PriceSource
	SourceState SourceStateStorage
	Now         func() time.Time
}

func (c DefaultSavedSkinPriceCollector) Collect(ctx context.Context, saved SavedSkin, params UpdateSavedSkinPriceParams) (UpdateSavedSkinPriceResult, error) {
	currency := NormalizeCurrencyCode(params.Currency)
	result := UpdateSavedSkinPriceResult{
		MarketHashName:    params.MarketHashName,
		SteamPageURL:      saved.SteamPageURL,
		SteamPriceText:    saved.SteamPriceText,
		SteamUpdatedAt:    saved.SteamUpdatedAt,
		LisSkinsPageURL:   saved.LisSkinsPageURL,
		LisSkinsPriceText: saved.LisSkinsPriceText,
		LisSkinsUpdatedAt: saved.LisSkinsUpdatedAt,
		CSTMPageURL:       saved.CSTMPageURL,
		CSTMPriceText:     saved.CSTMPriceText,
		CSTMUpdatedAt:     saved.CSTMUpdatedAt,
		Prices:            append([]PriceSnapshotView(nil), saved.Prices...),
		Currency:          currency,
	}

	updatedAny := false
	var firstErr error
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, source := range c.Sources {
		if source == nil {
			continue
		}
		source := source
		wg.Add(1)
		go func() {
			defer wg.Done()
			quote, fetchErr := source.FetchPrice(ctx, params.MarketHashName, currency)
			now := quote.FetchedAt
			if now.IsZero() {
				now = c.now()
			}
			mu.Lock()
			defer mu.Unlock()
			if fetchErr != nil {
				classified := classifySourceError("skins.update_one."+source.ID(), source.ID(), fetchErr)
				if firstErr == nil {
					firstErr = classified
				}
				c.recordSourceError(ctx, source.ID(), classified.Error(), c.now())
				return
			}
			if quote.Source == "" {
				quote.Source = source.ID()
			}
			if quote.SourceLabel == "" {
				quote.SourceLabel = source.Label()
			}
			if quote.Currency == "" {
				quote.Currency = currency
			}
			if quote.PriceText == "" && quote.PriceCents == nil {
				classified := errx.E("skins.update_one."+source.ID(), errx.CodeExternal, source.ID()+" returned empty price", nil)
				if firstErr == nil {
					firstErr = classified
				}
				c.recordSourceError(ctx, source.ID(), classified.Error(), c.now())
				return
			}
			updatedAny = true
			snapshot := PriceSnapshotView{
				Source:      quote.Source,
				SourceLabel: quote.SourceLabel,
				PageURL:     quote.PageURL,
				PriceText:   quote.PriceText,
				PriceCents:  quote.PriceCents,
				Currency:    NormalizeCurrencyCode(quote.Currency),
				FetchedAt:   now,
				Status:      "ok",
			}
			result.Prices = upsertSnapshotView(result.Prices, snapshot)
			applyLegacySourcePrice(&result, snapshot)
			c.recordSourceSuccess(ctx, source.ID(), now)
		}()
	}

	wg.Wait()

	if !updatedAny {
		if firstErr != nil {
			return UpdateSavedSkinPriceResult{}, firstErr
		}
		return UpdateSavedSkinPriceResult{}, errx.E("skins.update_one.no_sources", errx.CodeInternal, "no price sources configured", nil)
	}
	return result, nil
}

func (c DefaultSavedSkinPriceCollector) recordSourceSuccess(ctx context.Context, source string, at time.Time) {
	if c.SourceState != nil {
		_ = c.SourceState.RecordSourceSuccess(ctx, source, at)
	}
}

func (c DefaultSavedSkinPriceCollector) recordSourceError(ctx context.Context, source, message string, at time.Time) {
	if c.SourceState != nil {
		_ = c.SourceState.RecordSourceError(ctx, source, message, at)
	}
}

func upsertSnapshotView(items []PriceSnapshotView, next PriceSnapshotView) []PriceSnapshotView {
	for i, item := range items {
		if item.Source == next.Source {
			items[i] = next
			return items
		}
	}
	return append(items, next)
}

func applyLegacySourcePrice(result *UpdateSavedSkinPriceResult, snapshot PriceSnapshotView) {
	switch snapshot.Source {
	case "steam":
		result.SteamPageURL = snapshot.PageURL
		result.SteamPriceText = snapshot.PriceText
		result.SteamUpdatedAt = snapshot.FetchedAt
	case "lisskins":
		result.LisSkinsPageURL = snapshot.PageURL
		result.LisSkinsPriceText = snapshot.PriceText
		result.LisSkinsUpdatedAt = snapshot.FetchedAt
	case "cstm":
		result.CSTMPageURL = snapshot.PageURL
		result.CSTMPriceText = snapshot.PriceText
		result.CSTMUpdatedAt = snapshot.FetchedAt
	}
}

func (c DefaultSavedSkinPriceCollector) now() time.Time {
	if c.Now != nil {
		return c.Now().UTC()
	}
	return time.Now().UTC()
}

func classifySourceError(op, source string, err error) error {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "context deadline exceeded"):
		return errx.E(op, errx.CodeTimeout, source+" request timed out", err)
	case strings.Contains(message, "bad status"), strings.Contains(message, "unsuccessful"):
		return errx.E(op, errx.CodeUnavailable, source+" is temporarily unavailable", err)
	case strings.Contains(message, "not found"):
		return errx.E(op, errx.CodeNotFound, "skin not found on "+source, err)
	case strings.Contains(message, "decode"):
		return errx.E(op, errx.CodeExternal, "failed to decode "+source+" response", err)
	default:
		return errx.E(op, errx.CodeExternal, source+" request failed", err)
	}
}
