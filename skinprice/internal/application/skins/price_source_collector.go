package skins

import (
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"strings"
	"sync"
	"time"
)

type DefaultSavedSkinPriceCollector struct {
	SteamSource    MarketPriceReader
	LisSkinsSource MarketPriceReader
	CSTMSource     MarketPriceReader
	LisSkinsPages  MarketPageURLBuilder
	CSTMPages      MarketPageURLBuilder
	Now            func() time.Time
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
		Currency:          currency,
	}

	updatedAny := false
	steamErr := error(nil)
	lisSkinsErr := error(nil)
	cstmErr := error(nil)
	var mu sync.Mutex
	var wg sync.WaitGroup

	if c.SteamSource != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			price, fetchErr := c.SteamSource.GetByMarketHashName(ctx, params.MarketHashName, currency)
			mu.Lock()
			defer mu.Unlock()
			if fetchErr != nil {
				steamErr = classifySourceError("skins.update_one.steam", "steam", fetchErr)
				return
			}
			if price == nil {
				steamErr = errx.E("skins.update_one.steam", errx.CodeExternal, "steam returned empty price", nil)
				return
			}
			updatedAny = true
			result.SteamPageURL = FirstNonEmpty(price.PageURL, result.SteamPageURL)
			result.SteamPriceText = NormalizePriceText(price, currency)
			result.SteamUpdatedAt = c.now()
		}()
	}

	if c.LisSkinsSource != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			price, fetchErr := c.LisSkinsSource.GetByMarketHashName(ctx, params.MarketHashName, LisSkinsCurrency)
			mu.Lock()
			defer mu.Unlock()
			if fetchErr != nil {
				lisSkinsErr = classifySourceError("skins.update_one.lisskins", "lisskins", fetchErr)
				return
			}
			if price == nil {
				lisSkinsErr = errx.E("skins.update_one.lisskins", errx.CodeExternal, "lisskins returned empty price", nil)
				return
			}
			updatedAny = true
			fallbackURL := ""
			if c.LisSkinsPages != nil {
				fallbackURL = c.LisSkinsPages.BuildMarketPageURL(params.MarketHashName)
			}
			result.LisSkinsPageURL = FirstNonEmpty(price.PageURL, result.LisSkinsPageURL, fallbackURL)
			result.LisSkinsPriceText = NormalizePriceText(price, LisSkinsCurrency)
			result.LisSkinsUpdatedAt = c.now()
		}()
	}

	if c.CSTMSource != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			price, fetchErr := c.CSTMSource.GetByMarketHashName(ctx, params.MarketHashName, currency)
			mu.Lock()
			defer mu.Unlock()
			if fetchErr != nil {
				cstmErr = classifySourceError("skins.update_one.cstm", "cstm", fetchErr)
				return
			}
			if price == nil {
				cstmErr = errx.E("skins.update_one.cstm", errx.CodeExternal, "cstm returned empty price", nil)
				return
			}
			updatedAny = true
			fallbackURL := ""
			if c.CSTMPages != nil {
				fallbackURL = c.CSTMPages.BuildMarketPageURL(params.MarketHashName)
			}
			result.CSTMPageURL = FirstNonEmpty(price.PageURL, result.CSTMPageURL, fallbackURL)
			result.CSTMPriceText = NormalizePriceText(price, currency)
			result.CSTMUpdatedAt = c.now()
		}()
	}

	wg.Wait()

	if !updatedAny {
		if steamErr != nil {
			return UpdateSavedSkinPriceResult{}, steamErr
		}
		if lisSkinsErr != nil {
			return UpdateSavedSkinPriceResult{}, lisSkinsErr
		}
		if cstmErr != nil {
			return UpdateSavedSkinPriceResult{}, cstmErr
		}
		return UpdateSavedSkinPriceResult{}, errx.E("skins.update_one.no_sources", errx.CodeInternal, "no price sources configured", nil)
	}
	return result, nil
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
