package csgo

import (
	"time"

	"github.com/google/uuid"
)

type MarketPriceType string

const (
	MarketPriceTypeSellNow     MarketPriceType = "sell_now"
	MarketPriceTypeLastSale    MarketPriceType = "last_sale"
	MarketPriceTypeRecommended MarketPriceType = "recommended"
	MarketPriceTypeAverage     MarketPriceType = "average"
)

type MarketPrice struct {
	ID              uuid.UUID
	SkinID          uuid.UUID
	Platform        Platform
	Price           float64
	PriceType       MarketPriceType
	Currency        string
	AvailableVolume int
	FetchedAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewMarketPrice(
	SkinID uuid.UUID,
	Platform Platform,
	Price float64,
	PriceType MarketPriceType,
	Currency string,
	AvailableVolume int,
	FetchedAt time.Time,
) *MarketPrice {
	now := time.Now()
	return &MarketPrice{
		ID:              uuid.New(),
		SkinID:          SkinID,
		Platform:        Platform,
		Price:           Price,
		PriceType:       PriceType,
		Currency:        Currency,
		AvailableVolume: AvailableVolume,
		FetchedAt:       FetchedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
