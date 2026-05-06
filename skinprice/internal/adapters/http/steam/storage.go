package steam

import (
	"SkinPrice/skinprice/internal/application"
	"SkinPrice/skinprice/internal/application/skins"
	"net/http"
	"net/url"
	"strconv"
)

type Storage struct {
	Client  *http.Client
	BaseUrl string
}

func (s *Storage) GetList(criteria skins.SearchCriteria, params *application.Pagination) (skins.NewSkinsList, error) {
	q := buildSteamMarketSearchParams(params)
	if criteria.MarketHashName != nil {
		q.Set("query", *criteria.MarketHashName)
	}

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
