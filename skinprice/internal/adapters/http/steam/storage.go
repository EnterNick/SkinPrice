package steam

import (
	"SkinPrice/skinprice/internal/application"
	"SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/shared/logx"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Storage struct {
	Client  *http.Client
	BaseURL string
	Logger  *slog.Logger
}

type steamMarketSearchResponse struct {
	Success      bool                    `json:"success"`
	Start        int                     `json:"start"`
	Pagesize     int                     `json:"pagesize"`
	TotalCount   int                     `json:"total_count"`
	Results      []steamMarketSearchItem `json:"results"`
	ErrorMessage string                  `json:"error"`
}

type steamMarketPriceOverviewResponse struct {
	Success      bool   `json:"success"`
	LowestPrice  string `json:"lowest_price"`
	MedianPrice  string `json:"median_price"`
	Volume       string `json:"volume"`
	ErrorMessage string `json:"error"`
}

type steamMarketSearchItem struct {
	Name          string `json:"name"`
	HashName      string `json:"hash_name"`
	SellListings  int64  `json:"sell_listings"`
	SellPrice     *int64 `json:"sell_price"`
	SellPriceText string `json:"sell_price_text"`
	AssetDesc     struct {
		IconURL   string `json:"icon_url"`
		NameColor string `json:"name_color"`
	} `json:"asset_description"`
}

const requestTimeout = 8 * time.Second

func (s *Storage) GetList(criteria skins.SearchCriteria, params *application.Pagination) (_ skins.NewSkinsList, err error) {
	logger := logx.WithComponent(s.Logger, "steam_storage")
	q := buildSteamMarketSearchParams(params)
	if criteria.MarketHashName != nil {
		q.Set("query", *criteria.MarketHashName)
	}

	endpoint := fmt.Sprintf("%s/search/render/?%s", s.BaseURL, q.Encode())
	startedAt := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return skins.NewSkinsList{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	setSteamHeaders(req)

	resp, err := s.Client.Do(req)
	if err != nil {
		logger.Error("steam search request failed",
			append([]any{
				slog.String("operation", "search"),
				slog.String("source", "steam"),
				slog.Duration("duration", time.Since(startedAt)),
			}, logx.ErrAttrs(err)...)...,
		)
		return skins.NewSkinsList{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close steam list response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logger.Warn("steam search returned non-200",
			slog.String("operation", "search"),
			slog.String("source", "steam"),
			slog.Int("status_code", resp.StatusCode),
			slog.Duration("duration", time.Since(startedAt)),
		)
		return skins.NewSkinsList{}, fmt.Errorf("%w: %d", skins.ErrNewSkinsRequestBadStatus, resp.StatusCode)
	}

	var payload steamMarketSearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		logger.Error("steam search decode failed",
			append([]any{
				slog.String("operation", "search"),
				slog.String("source", "steam"),
				slog.Duration("duration", time.Since(startedAt)),
			}, logx.ErrAttrs(err)...)...,
		)
		return skins.NewSkinsList{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}
	if !payload.Success {
		logger.Warn("steam search unsuccessful",
			slog.String("operation", "search"),
			slog.String("source", "steam"),
			slog.String("message", payload.ErrorMessage),
			slog.Duration("duration", time.Since(startedAt)),
		)
		return skins.NewSkinsList{}, fmt.Errorf("%w: %s", skins.ErrNewSkinsResponseUnsuccess, payload.ErrorMessage)
	}

	items := make([]skins.NewSkin, 0, len(payload.Results))
	for _, item := range payload.Results {
		items = append(items, skins.NewSkin{
			MarketHashName: item.HashName,
			DisplayName:    item.Name,
			NameColor:      item.AssetDesc.NameColor,
			SellListings:   item.SellListings,
			PriceCents:     item.SellPrice,
			PriceText:      item.SellPriceText,
			IconURL:        item.AssetDesc.IconURL,
			PageURL:        fmt.Sprintf("%s/listings/730/%s", s.BaseURL, url.PathEscape(item.HashName)),
		})
	}

	logger.Debug("steam search completed",
		slog.String("operation", "search"),
		slog.String("source", "steam"),
		slog.Int("items", len(items)),
		slog.Int("total_count", payload.TotalCount),
		slog.Duration("duration", time.Since(startedAt)),
	)

	return skins.NewSkinsList{
		Items:      items,
		TotalCount: payload.TotalCount,
		Offset:     payload.Start,
		Limit:      payload.Pagesize,
	}, nil
}

func buildSteamMarketSearchParams(params *application.Pagination) url.Values {
	q := url.Values{}

	q.Set("start", strconv.Itoa(params.Offset))
	q.Set("count", strconv.Itoa(params.Limit))
	q.Set("search_descriptions", "0")
	q.Set("appid", "730")
	q.Set("norender", "1")

	return q
}

func (s *Storage) GetByMarketHashName(marketHashName, currency string) (_ *skins.NewSkin, err error) {
	logger := logx.WithComponent(s.Logger, "steam_storage")
	q := url.Values{}
	q.Set("appid", "730")
	q.Set("market_hash_name", marketHashName)
	if currency != "" {
		q.Set("currency", currency)
	}

	endpoint := fmt.Sprintf("%s/priceoverview/?%s", s.BaseURL, q.Encode())
	startedAt := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	setSteamHeaders(req)

	resp, err := s.Client.Do(req)
	if err != nil {
		logger.Error("steam price request failed",
			append([]any{
				slog.String("operation", "price_overview"),
				slog.String("source", "steam"),
				slog.String("market_hash_name", marketHashName),
				slog.String("currency", currency),
				slog.Duration("duration", time.Since(startedAt)),
			}, logx.ErrAttrs(err)...)...,
		)
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close steam item response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logger.Warn("steam price returned non-200",
			slog.String("operation", "price_overview"),
			slog.String("source", "steam"),
			slog.String("market_hash_name", marketHashName),
			slog.Int("status_code", resp.StatusCode),
			slog.Duration("duration", time.Since(startedAt)),
		)
		return nil, fmt.Errorf("%w: %d", skins.ErrNewSkinsRequestBadStatus, resp.StatusCode)
	}

	var payload steamMarketPriceOverviewResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		logger.Error("steam price decode failed",
			append([]any{
				slog.String("operation", "price_overview"),
				slog.String("source", "steam"),
				slog.String("market_hash_name", marketHashName),
				slog.Duration("duration", time.Since(startedAt)),
			}, logx.ErrAttrs(err)...)...,
		)
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}
	if !payload.Success {
		logger.Warn("steam price unsuccessful",
			slog.String("operation", "price_overview"),
			slog.String("source", "steam"),
			slog.String("market_hash_name", marketHashName),
			slog.String("message", payload.ErrorMessage),
			slog.Duration("duration", time.Since(startedAt)),
		)
		return nil, fmt.Errorf("%w: %s", skins.ErrNewSkinsResponseUnsuccess, payload.ErrorMessage)
	}
	priceText := payload.LowestPrice
	if priceText == "" {
		priceText = payload.MedianPrice
	}
	if priceText == "" {
		logger.Warn("steam price response missing price",
			slog.String("operation", "price_overview"),
			slog.String("source", "steam"),
			slog.String("market_hash_name", marketHashName),
			slog.Duration("duration", time.Since(startedAt)),
		)
		return nil, fmt.Errorf("%w: missing price overview", skins.ErrNewSkinsResponseUnsuccess)
	}

	result := &skins.NewSkin{
		MarketHashName: marketHashName,
		DisplayName:    marketHashName,
		PriceText:      priceText,
		PageURL:        fmt.Sprintf("%s/listings/730/%s", s.BaseURL, url.PathEscape(marketHashName)),
	}
	logger.Debug("steam price completed",
		slog.String("operation", "price_overview"),
		slog.String("source", "steam"),
		slog.String("market_hash_name", marketHashName),
		slog.String("currency", currency),
		slog.Duration("duration", time.Since(startedAt)),
	)
	return result, nil
}

func setSteamHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://steamcommunity.com/market/search?appid=730")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
}
