package csgo

import (
	"time"

	"github.com/google/uuid"
)

type Skin struct {
	ID                uuid.UUID
	MarketHashName    string
	Name              string
	SkinName          string
	Exterior          string
	Rarity            string
	Collection        string
	ImageURL          string
	SteamMarketURL    string
	LisSkinsMarketURL string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewSkin(
	MarketHashName string,
	Name string,
	SkinName string,
	Exterior string,
	Rarity string,
	Collection string,
	ImageURL string,
	SteamMarketURL string,
	LisSkinsMarketURL string,
) *Skin {
	now := time.Now()
	return &Skin{
		ID:                uuid.New(),
		MarketHashName:    MarketHashName,
		Name:              Name,
		SkinName:          SkinName,
		Exterior:          Exterior,
		Rarity:            Rarity,
		Collection:        Collection,
		ImageURL:          ImageURL,
		SteamMarketURL:    SteamMarketURL,
		LisSkinsMarketURL: LisSkinsMarketURL,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
