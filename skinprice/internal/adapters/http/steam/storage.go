package steam

import (
	"SkinPrice/skinprice/internal/application"
	"SkinPrice/skinprice/internal/application/skins"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Storage struct {
	Client  *http.Client
	BaseURL string
}

type steamMarketSearchResponse struct {
	Success      bool                    `json:"success"`
	Start        int                     `json:"start"`
	Pagesize     int                     `json:"pagesize"`
	TotalCount   int                     `json:"total_count"`
	Results      []steamMarketSearchItem `json:"results"`
	ErrorMessage string                  `json:"error"`
}

type steamMarketSearchItem struct {
	Name          string `json:"name"`
	HashName      string `json:"hash_name"`
	SellListings  int64  `json:"sell_listings"`
	SellPrice     *int64 `json:"sell_price"`
	SellPriceText string `json:"sell_price_text"`
	AssetDesc     struct {
		IconURL string `json:"icon_url"`
	} `json:"asset_description"`
}

func (s *Storage) GetList(criteria skins.SearchCriteria, params *application.Pagination) (_ skins.NewSkinsList, err error) {
	q := buildSteamMarketSearchParams(params)
	if criteria.MarketHashName != nil {
		q.Set("query", *criteria.MarketHashName)
	}

	endpoint := fmt.Sprintf("%s/search/render/?%s", s.BaseURL, q.Encode())
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return skins.NewSkinsList{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return skins.NewSkinsList{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close steam list response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return skins.NewSkinsList{}, fmt.Errorf("%w: %d", skins.ErrNewSkinsRequestBadStatus, resp.StatusCode)
	}

	var payload steamMarketSearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return skins.NewSkinsList{}, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}
	if !payload.Success {
		return skins.NewSkinsList{}, fmt.Errorf("%w: %s", skins.ErrNewSkinsResponseUnsuccess, payload.ErrorMessage)
	}

	items := make([]skins.NewSkin, 0, len(payload.Results))
	for _, item := range payload.Results {
		items = append(items, skins.NewSkin{
			MarketHashName: item.HashName,
			DisplayName:    item.Name,
			SellListings:   item.SellListings,
			PriceCents:     item.SellPrice,
			PriceText:      item.SellPriceText,
			IconURL:        item.AssetDesc.IconURL,
			PageURL:        fmt.Sprintf("%s/listings/730/%s", s.BaseURL, url.PathEscape(item.HashName)),
		})
	}

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
	q := buildSteamMarketSearchParams(&application.Pagination{Limit: 1, Offset: 0})
	q.Set("query", marketHashName)
	if currency != "" {
		q.Set("currency", currency)
	}

	endpoint := fmt.Sprintf("%s/search/render/?%s", s.BaseURL, q.Encode())
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsRequestFailed, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close steam item response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", skins.ErrNewSkinsRequestBadStatus, resp.StatusCode)
	}

	var payload steamMarketSearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("%w: %w", skins.ErrNewSkinsResponseDecodeFail, err)
	}
	if !payload.Success || len(payload.Results) == 0 {
		return nil, fmt.Errorf("%w: %s", skins.ErrNewSkinsResponseUnsuccess, payload.ErrorMessage)
	}

	for _, item := range payload.Results {
		if item.HashName == marketHashName {
			return &skins.NewSkin{
				MarketHashName: item.HashName,
				DisplayName:    item.Name,
				SellListings:   item.SellListings,
				PriceCents:     item.SellPrice,
				PriceText:      item.SellPriceText,
				IconURL:        item.AssetDesc.IconURL,
				PageURL:        fmt.Sprintf("%s/listings/730/%s", s.BaseURL, url.PathEscape(item.HashName)),
			}, nil
		}
	}

	return nil, fmt.Errorf("%w: skin not found", skins.ErrNewSkinsResponseUnsuccess)
}
