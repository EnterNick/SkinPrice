package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"fmt"
)

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) Save(params appskins.SaveSkinParams) error {
	db := s.Conn.DB()

	_, err := db.Exec(
		`INSERT INTO skins (market_hash_name, display_name, icon_url, page_url) VALUES ($1, $2, $3, $4) ON CONFLICT (market_hash_name) DO NOTHING`,
		params.MarketHashName,
		params.DisplayName,
		params.IconURL,
		params.PageURL,
	)
	if err != nil {
		_, sqliteErr := db.Exec(
			`INSERT OR IGNORE INTO skins (market_hash_name, display_name, icon_url, page_url) VALUES (?, ?, ?, ?)`,
			params.MarketHashName,
			params.DisplayName,
			params.IconURL,
			params.PageURL,
		)
		if sqliteErr != nil {
			return fmt.Errorf("save skin: %w", err)
		}
	}

	return nil
}
