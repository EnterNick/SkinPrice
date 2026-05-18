package lisskins

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

type marketSearchResponse struct {
	Data       []marketSearchItem `json:"data"`
	Items      []marketSearchItem `json:"items"`
	Results    []marketSearchItem `json:"results"`
	TotalCount int                `json:"total_count"`
	Count      int                `json:"count"`
	NextCursor string             `json:"next_cursor"`
	Cursor     string             `json:"cursor"`
	Success    *bool              `json:"success"`
	Error      string             `json:"error"`
	Message    string             `json:"message"`
}

type marketSearchItem struct {
	Name          string `json:"name"`
	Title         string `json:"title"`
	HashName      string `json:"hash_name"`
	MarketHash    string `json:"market_hash_name"`
	Listings      int64  `json:"sell_listings"`
	Count         int64  `json:"count"`
	SellPrice     *int64 `json:"sell_price"`
	Price         *int64 `json:"price"`
	SellPriceText string `json:"sell_price_text"`
	PriceText     string `json:"price_text"`
	IconURL       string `json:"icon_url"`
	Image         string `json:"image"`
	AssetDesc     struct {
		IconURL string `json:"icon_url"`
	} `json:"asset_description"`
}

const requestTimeout = 8 * time.Second

func (s *Storage) GetList(criteria skins.SearchCriteria, params *application.Pagination) (skins.NewSkinsList, error) {
	logger := logx.WithComponent(s.Logger, "lisskins_storage")
	startedAt := time.Now()
	q := buildLisSkinsMarketSearchParams(criteria, params, "")
	endpoint := fmt.Sprintf("%s/market/search?%s", s.BaseURL, q.Encode())
	payload, err := s.fetch(endpoint)
	if err != nil {
		logger.Error("lisskins search failed",
			append([]any{
				slog.String("operation", "search"),
				slog.String("source", "lisskins"),
				slog.Duration("duration", time.Since(startedAt)),
			}, logx.ErrAttrs(err)...)...,
		)
		return skins.NewSkinsList{}, err
	}

	items := payload.extractItems()
	result := make([]skins.NewSkin, 0, len(items))
	for _, item := range items {
		result = append(result, mapItem(item, s.BaseURL))
	}

	totalCount := payload.TotalCount
	if totalCount == 0 {
		totalCount = payload.Count
	}

	logger.Debug("lisskins search completed",
		slog.String("operation", "search"),
		slog.String("source", "lisskins"),
		slog.Int("items", len(result)),
		slog.Int("total_count", totalCount),
		slog.Duration("duration", time.Since(startedAt)),
	)

	return skins.NewSkinsList{
		Items:      result,
		TotalCount: totalCount,
		Offset:     params.Offset,
		Limit:      params.Limit,
	}, nil
}

func (s *Storage) GetByMarketHashName(marketHashName, currency string) (*skins.NewSkin, error) {
	logger := logx.WithComponent(s.Logger, "lisskins_storage")
	startedAt := time.Now()
	q := buildLisSkinsMarketSearchParams(skins.SearchCriteria{MarketHashName: &marketHashName}, &application.Pagination{Limit: 20, Offset: 0}, currency)
	endpoint := fmt.Sprintf("%s/market/search?%s", s.BaseURL, q.Encode())
	payload, err := s.fetch(endpoint)
	if err != nil {
		logger.Error("lisskins price search failed",
			append([]any{
				slog.String("operation", "lookup"),
				slog.String("source", "lisskins"),
				slog.String("market_hash_name", marketHashName),
				slog.String("currency", currency),
				slog.Duration("duration", time.Since(startedAt)),
			}, logx.ErrAttrs(err)...)...,
		)
		return nil, err
	}

	for _, item := range payload.extractItems() {
		hash := firstNonEmpty(item.HashName, item.MarketHash)
		if hash == marketHashName {
			skin := mapItem(item, s.BaseURL)
			logger.Debug("lisskins lookup completed",
				slog.String("operation", "lookup"),
				slog.String("source", "lisskins"),
				slog.String("market_hash_name", marketHashName),
				slog.String("currency", currency),
				slog.Duration("duration", time.Since(startedAt)),
			)
			return &skin, nil
		}
	}
	logger.Warn("lisskins lookup returned no result",
		slog.String("operation", "lookup"),
		slog.String("source", "lisskins"),
		slog.String("market_hash_name", marketHashName),
		slog.String("currency", currency),
		slog.Duration("duration", time.Since(startedAt)),
	)
	return nil, fmt.Errorf("%w: skin not found", skins.ErrNewSkinsResponseUnsuccess)
}

func (s *Storage) fetch(endpoint string) (_ marketSearchResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return marketSearchResponse{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	setLisSkinsHeaders(req)

	resp, err := s.Client.Do(req)
	if err != nil {
		return marketSearchResponse{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close lisskins response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return marketSearchResponse{}, fmt.Errorf("%w: %d", skins.ErrNewSkinsRequestBadStatus, resp.StatusCode)
	}

	var payload marketSearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return marketSearchResponse{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}
	if payload.Success != nil && !*payload.Success {
		return marketSearchResponse{}, fmt.Errorf("%w: %s", skins.ErrNewSkinsResponseUnsuccess, firstNonEmpty(payload.Error, payload.Message))
	}
	return payload, nil
}

func setLisSkinsHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://lis-skins.ru/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
}

func buildLisSkinsMarketSearchParams(criteria skins.SearchCriteria, params *application.Pagination, currency string) url.Values {
	q := url.Values{}
	q.Set("app_id", "730")
	q.Set("limit", strconv.Itoa(params.Limit))
	if params.Offset > 0 {
		q.Set("cursor", strconv.Itoa(params.Offset))
	}
	if criteria.MarketHashName != nil && *criteria.MarketHashName != "" {
		q.Set("name", *criteria.MarketHashName)
	}
	if currency != "" {
		q.Set("currency", currency)
	}
	return q
}

func (r marketSearchResponse) extractItems() []marketSearchItem {
	if len(r.Data) > 0 {
		return r.Data
	}
	if len(r.Items) > 0 {
		return r.Items
	}
	return r.Results
}

func mapItem(item marketSearchItem, baseURL string) skins.NewSkin {
	hash := firstNonEmpty(item.HashName, item.MarketHash)
	icon := firstNonEmpty(item.IconURL, item.Image, item.AssetDesc.IconURL)
	priceText := firstNonEmpty(item.SellPriceText, item.PriceText)
	sellListings := item.Listings
	if sellListings == 0 {
		sellListings = item.Count
	}
	priceCents := item.SellPrice
	if priceCents == nil {
		priceCents = item.Price
	}
	return skins.NewSkin{
		MarketHashName: hash,
		DisplayName:    firstNonEmpty(item.Name, item.Title, hash),
		SellListings:   sellListings,
		PriceCents:     priceCents,
		PriceText:      priceText,
		IconURL:        icon,
		PageURL:        fmt.Sprintf("%s/market/csgo/%s", baseURL, url.PathEscape(hash)),
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
