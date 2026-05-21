package cstm

import (
	"SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/shared/logx"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Storage struct {
	Client         *http.Client
	BaseURL        string
	RequestTimeout time.Duration
	Logger         *slog.Logger
	CacheTTL       time.Duration

	mu    sync.RWMutex
	cache map[string]cachedPrices
}

type cachedPrices struct {
	items     map[string]priceListItem
	expiresAt time.Time
}

type priceListItem struct {
	MarketHashName string `json:"market_hash_name"`
	Volume         string `json:"volume"`
	Price          string `json:"price"`
}

type priceListEnvelope struct {
	Success *bool           `json:"success"`
	Error   string          `json:"error"`
	Message string          `json:"message"`
	Items   []priceListItem `json:"items"`
	Data    []priceListItem `json:"data"`
	Result  []priceListItem `json:"result"`
	Prices  []priceListItem `json:"prices"`
}

func (s *Storage) GetByMarketHashName(ctx context.Context, marketHashName, currency string) (*skins.NewSkin, error) {
	logger := logx.WithComponent(s.Logger, "cstm_storage")
	startedAt := time.Now()
	code := currencyToPriceListCode(currency)
	items, err := s.getPriceList(ctx, code)
	if err != nil {
		logger.Error("cstm price lookup failed",
			append([]any{
				slog.String("operation", "lookup"),
				slog.String("source", "cstm"),
				slog.String("market_hash_name", marketHashName),
				slog.String("currency", code),
				slog.Duration("duration", time.Since(startedAt)),
			}, logx.ErrAttrs(err)...)...,
		)
		return nil, err
	}

	item, ok := items[marketHashName]
	if !ok {
		logger.Warn("cstm lookup returned no result",
			slog.String("operation", "lookup"),
			slog.String("source", "cstm"),
			slog.String("market_hash_name", marketHashName),
			slog.String("currency", code),
			slog.Duration("duration", time.Since(startedAt)),
		)
		return nil, fmt.Errorf("%w: skin not found", skins.ErrNewSkinsResponseUnsuccess)
	}

	priceCents, err := parsePriceToCents(item.Price)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}
	skin := &skins.NewSkin{
		MarketHashName: marketHashName,
		DisplayName:    marketHashName,
		PageURL:        BuildMarketPageURL(s.BaseURL, marketHashName),
		PriceCents:     &priceCents,
		PriceText:      formatRawPriceText(item.Price, code),
	}
	logger.Debug("cstm lookup completed",
		slog.String("operation", "lookup"),
		slog.String("source", "cstm"),
		slog.String("market_hash_name", marketHashName),
		slog.String("currency", code),
		slog.Duration("duration", time.Since(startedAt)),
	)
	return skin, nil
}

func (s *Storage) getPriceList(ctx context.Context, currencyCode string) (map[string]priceListItem, error) {
	now := time.Now()
	s.mu.RLock()
	if entry, ok := s.cache[currencyCode]; ok && now.Before(entry.expiresAt) {
		items := entry.items
		s.mu.RUnlock()
		return items, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cache == nil {
		s.cache = make(map[string]cachedPrices)
	}
	if entry, ok := s.cache[currencyCode]; ok && now.Before(entry.expiresAt) {
		return entry.items, nil
	}

	items, err := s.fetchPriceList(ctx, currencyCode)
	if err != nil {
		return nil, err
	}
	s.cache[currencyCode] = cachedPrices{
		items:     items,
		expiresAt: now.Add(cacheTTLOrDefault(s.CacheTTL)),
	}
	return items, nil
}

func (s *Storage) fetchPriceList(ctx context.Context, currencyCode string) (map[string]priceListItem, error) {
	endpoint := fmt.Sprintf("%s/api/v2/prices/%s.json", strings.TrimRight(s.BaseURL, "/"), currencyCode)
	ctx, cancel := context.WithTimeout(ctx, requestTimeoutOrDefault(s.RequestTimeout))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "SkinPrice/1.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", skins.ErrNewSkinsRequestBadStatus, resp.StatusCode)
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}

	payload, err := decodePriceList(rawBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}

	items := make(map[string]priceListItem, len(payload))
	for _, item := range payload {
		if strings.TrimSpace(item.MarketHashName) == "" {
			continue
		}
		items[item.MarketHashName] = item
	}
	return items, nil
}

func decodePriceList(rawBody []byte) ([]priceListItem, error) {
	trimmed := bytes.TrimSpace(rawBody)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("empty response body")
	}

	switch trimmed[0] {
	case '[':
		var payload []priceListItem
		if err := json.Unmarshal(trimmed, &payload); err != nil {
			return nil, err
		}
		return payload, nil
	case '{':
		var envelope priceListEnvelope
		if err := json.Unmarshal(trimmed, &envelope); err != nil {
			return nil, err
		}
		if envelope.Success != nil && !*envelope.Success {
			return nil, fmt.Errorf("unsuccessful response: %s", firstNonEmpty(envelope.Error, envelope.Message))
		}
		switch {
		case len(envelope.Items) > 0:
			return envelope.Items, nil
		case len(envelope.Data) > 0:
			return envelope.Data, nil
		case len(envelope.Result) > 0:
			return envelope.Result, nil
		case len(envelope.Prices) > 0:
			return envelope.Prices, nil
		default:
			return nil, fmt.Errorf("unsupported response object shape")
		}
	default:
		return nil, fmt.Errorf("unexpected response prefix %q", string(trimmed[:1]))
	}
}

func BuildMarketPageURL(baseURL, marketHashName string) string {
	return strings.TrimRight(baseURL, "/") + "/en/" + url.PathEscape(marketHashName)
}

func (s *Storage) BuildMarketPageURL(marketHashName string) string {
	return BuildMarketPageURL(s.BaseURL, marketHashName)
}

func currencyToPriceListCode(currency string) string {
	switch strings.TrimSpace(strings.ToUpper(currency)) {
	case "1", "USD":
		return "USD"
	case "3", "EUR":
		return "EUR"
	case "5", "RUB":
		return "RUB"
	default:
		return "USD"
	}
}

func parsePriceToCents(value string) (int64, error) {
	parts := strings.SplitN(strings.TrimSpace(value), ".", 3)
	if len(parts) == 0 || parts[0] == "" {
		return 0, fmt.Errorf("invalid price %q", value)
	}
	whole := int64(0)
	for _, ch := range parts[0] {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("invalid whole price %q", value)
		}
		whole = whole*10 + int64(ch-'0')
	}

	fraction := int64(0)
	if len(parts) == 2 {
		raw := parts[1]
		if len(raw) == 1 {
			raw += "0"
		}
		if len(raw) > 2 {
			raw = raw[:2]
		}
		for _, ch := range raw {
			if ch < '0' || ch > '9' {
				return 0, fmt.Errorf("invalid fraction price %q", value)
			}
			fraction = fraction*10 + int64(ch-'0')
		}
	}
	return whole*100 + fraction, nil
}

func formatRawPriceText(value, currencyCode string) string {
	switch currencyCode {
	case "EUR":
		return "€" + value
	case "RUB":
		return value + " ₽"
	default:
		return "$" + value
	}
}

func cacheTTLOrDefault(value time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	return 5 * time.Minute
}

func requestTimeoutOrDefault(value time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	return 15 * time.Second
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
