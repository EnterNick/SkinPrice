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
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	Client        *http.Client
	BaseURL       string
	Logger        *slog.Logger
	TokenProvider LisSkinsTokenProvider
}

type LisSkinsTokenProvider interface {
	Execute() (string, error)
}

const lisSkinsAuthHeader = "Authorization"
const lisSkinsMarketBaseURL = "https://lis-skins.com"

func buildLisSkinsAuthHeaderValue(token string) string {
	return "Bearer " + token
}

type marketSearchResponse struct {
	Data       []marketSearchItem `json:"data"`
	Items      []marketSearchItem `json:"items"`
	Results    []marketSearchItem `json:"results"`
	TotalCount int                `json:"total_count"`
	Count      int                `json:"count"`
	NextCursor string             `json:"next_cursor"`
	Cursor     string             `json:"cursor"`
	Meta       struct {
		TotalCount int    `json:"total_count"`
		Count      int    `json:"count"`
		NextCursor string `json:"next_cursor"`
		Cursor     string `json:"cursor"`
	} `json:"meta"`
	Success *bool  `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type marketSearchItem struct {
	Name          string      `json:"name"`
	Title         string      `json:"title"`
	HashName      string      `json:"hash_name"`
	MarketHash    string      `json:"market_hash_name"`
	Listings      int64       `json:"sell_listings"`
	Count         int64       `json:"count"`
	SellPrice     json.Number `json:"sell_price"`
	Price         json.Number `json:"price"`
	SellPriceText string      `json:"sell_price_text"`
	PriceText     string      `json:"price_text"`
	IconURL       string      `json:"icon_url"`
	Image         string      `json:"image"`
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
		result = append(result, mapItem(item))
	}

	totalCount := payload.total()

	logger.Debug("lisskins search completed",
		slog.String("operation", "search"),
		slog.String("source", "lisskins"),
		slog.Int("items", len(result)),
		slog.Int("total_count", totalCount),
		slog.String("next_cursor", payload.nextCursor()),
		slog.Duration("duration", time.Since(startedAt)),
	)

	return skins.NewSkinsList{
		Items:      result,
		TotalCount: totalCount,
		Offset:     params.Offset,
		Limit:      params.Limit,
		NextCursor: payload.nextCursor(),
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
		if matchesMarketHashName(marketHashName, item) {
			skin := mapItem(item)
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

func matchesMarketHashName(expected string, item marketSearchItem) bool {
	expectedSlug := buildLisSkinsItemSlug(expected)
	candidates := []string{
		item.HashName,
		item.MarketHash,
		item.Name,
		item.Title,
		firstNonEmpty(item.HashName, item.MarketHash, item.Name, item.Title),
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if candidate == expected {
			return true
		}
		if buildLisSkinsItemSlug(candidate) == expectedSlug {
			return true
		}
	}
	return false
}

func (s *Storage) fetch(endpoint string) (_ marketSearchResponse, err error) {
	token := ""
	if s.TokenProvider != nil {
		token, err = s.TokenProvider.Execute()
		if err != nil {
			return marketSearchResponse{}, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return marketSearchResponse{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	setLisSkinsHeaders(req, token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return marketSearchResponse{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close lisskins response body: %w", closeErr))
		}
	}()

	if (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) && token != "" {
		return marketSearchResponse{}, skins.ErrLisSkinsTokenInvalid
	}
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

func setLisSkinsHeaders(req *http.Request, token string) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", lisSkinsMarketBaseURL+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	if token != "" {
		req.Header.Set(lisSkinsAuthHeader, buildLisSkinsAuthHeaderValue(token))
	}
}

func buildLisSkinsMarketSearchParams(criteria skins.SearchCriteria, params *application.Pagination, currency string) url.Values {
	q := url.Values{}
	q.Set("game", "csgo")
	q.Set("limit", strconv.Itoa(params.Limit))
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if criteria.MarketHashName != nil && *criteria.MarketHashName != "" {
		q.Add("names[]", *criteria.MarketHashName)
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

func (r marketSearchResponse) nextCursor() string {
	return firstNonEmpty(r.Meta.NextCursor, r.NextCursor)
}

func (r marketSearchResponse) total() int {
	switch {
	case r.TotalCount > 0:
		return r.TotalCount
	case r.Meta.TotalCount > 0:
		return r.Meta.TotalCount
	case r.Count > 0:
		return r.Count
	default:
		return r.Meta.Count
	}
}

func mapItem(item marketSearchItem) skins.NewSkin {
	hash := firstNonEmpty(item.HashName, item.MarketHash, item.Name, item.Title)
	icon := firstNonEmpty(item.IconURL, item.Image, item.AssetDesc.IconURL)
	priceValue := firstNonEmpty(item.SellPrice.String(), item.Price.String())
	priceCents := parseLisSkinsPriceCents(priceValue)
	priceText := firstNonEmpty(item.SellPriceText, item.PriceText, formatLisSkinsPriceText(priceCents))
	sellListings := item.Listings
	if sellListings == 0 {
		sellListings = item.Count
	}
	return skins.NewSkin{
		MarketHashName: hash,
		DisplayName:    firstNonEmpty(item.Name, item.Title, hash),
		SellListings:   sellListings,
		PriceCents:     priceCents,
		PriceText:      priceText,
		IconURL:        icon,
		PageURL:        BuildMarketPageURL(hash),
	}
}

func BuildMarketPageURL(name string) string {
	return fmt.Sprintf("%s/market/csgo/%s/", lisSkinsMarketBaseURL, buildLisSkinsItemSlug(name))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func parseLisSkinsPriceCents(value string) *int64 {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return nil
	}
	if strings.Contains(normalized, ".") {
		floatValue, err := strconv.ParseFloat(normalized, 64)
		if err != nil {
			return nil
		}
		cents := int64(math.Round(floatValue * 100))
		return &cents
	}
	cents, err := strconv.ParseInt(normalized, 10, 64)
	if err != nil {
		return nil
	}
	return &cents
}

func formatLisSkinsPriceText(priceCents *int64) string {
	if priceCents == nil {
		return ""
	}
	return fmt.Sprintf("$%.2f", float64(*priceCents)/100)
}

func buildLisSkinsItemSlug(name string) string {
	replacements := strings.NewReplacer(
		"StatTrak™", "StatTrak",
		"™", "",
		"★", "",
		"|", " ",
		"(", " ",
		")", " ",
		"[", " ",
		"]", " ",
		"{", " ",
		"}", " ",
		"'", "",
		"\"", "",
		",", " ",
		".", " ",
		"/", " ",
		"\\", " ",
		":", " ",
		";", " ",
		"+", " ",
	)
	normalized := strings.ToLower(strings.TrimSpace(replacements.Replace(name)))
	parts := strings.Fields(normalized)
	return strings.Join(parts, "-")
}
