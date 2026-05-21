package skins

import (
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"time"
)

type ReaderPriceSource struct {
	SourceID         string
	SourceLabel      string
	Reader           MarketPriceReader
	Pages            MarketPageURLBuilder
	CurrencyOverride string
	Now              func() time.Time
}

func (s ReaderPriceSource) ID() string {
	return s.SourceID
}

func (s ReaderPriceSource) Label() string {
	if s.SourceLabel != "" {
		return s.SourceLabel
	}
	return s.SourceID
}

func (s ReaderPriceSource) FetchPrice(ctx context.Context, marketHashName, currency string) (PriceQuote, error) {
	if s.Reader == nil {
		return PriceQuote{}, errx.E("skins.price_source.no_reader", errx.CodeInternal, "price source reader is not configured", nil)
	}
	sourceCurrency := NormalizeCurrencyCode(currency)
	if s.CurrencyOverride != "" {
		sourceCurrency = NormalizeCurrencyCode(s.CurrencyOverride)
	}
	price, err := s.Reader.GetByMarketHashName(ctx, marketHashName, sourceCurrency)
	if err != nil {
		return PriceQuote{}, err
	}
	if price == nil {
		return PriceQuote{}, errx.E("skins.price_source.empty_price", errx.CodeExternal, "price source returned empty price", nil)
	}
	pageURL := price.PageURL
	if pageURL == "" && s.Pages != nil {
		pageURL = s.Pages.BuildMarketPageURL(marketHashName)
	}
	return PriceQuote{
		Source:      s.ID(),
		SourceLabel: s.Label(),
		PageURL:     pageURL,
		PriceText:   NormalizePriceText(price, sourceCurrency),
		PriceCents:  price.PriceCents,
		Currency:    sourceCurrency,
		FetchedAt:   s.now(),
	}, nil
}

func (s ReaderPriceSource) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}
