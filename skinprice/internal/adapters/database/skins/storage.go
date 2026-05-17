package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"fmt"
)

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) Save(params appskins.SaveSkinParams) error {
	db := s.Conn.DB()

	if _, err := db.Exec(`ALTER TABLE skins ADD COLUMN market_hash_name TEXT`); err != nil {
		_ = err
	}
	if _, err := db.Exec(`ALTER TABLE skins ADD COLUMN display_name TEXT`); err != nil {
		_ = err
	}
	if _, err := db.Exec(`ALTER TABLE skins ADD COLUMN icon_url TEXT`); err != nil {
		_ = err
	}
	if _, err := db.Exec(`ALTER TABLE skins ADD COLUMN page_url TEXT`); err != nil {
		_ = err
	}
	if _, err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_skins_market_hash_name ON skins(market_hash_name)`); err != nil {
		return fmt.Errorf("create unique index: %w", err)
	}

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

func (s *Storage) GetSavedList(params *application.Pagination) (appskins.SavedSkinsList, error) {
	db := s.Conn.DB()

	var totalCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM skins`).Scan(&totalCount); err != nil {
		return appskins.SavedSkinsList{}, fmt.Errorf("count skins: %w", err)
	}

	rows, err := db.Query(
		`SELECT market_hash_name, display_name, icon_url, page_url FROM skins ORDER BY id DESC LIMIT $1 OFFSET $2`,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		rows, err = db.Query(
			`SELECT market_hash_name, display_name, icon_url, page_url FROM skins ORDER BY id DESC LIMIT ? OFFSET ?`,
			params.Limit,
			params.Offset,
		)
		if err != nil {
			return appskins.SavedSkinsList{}, fmt.Errorf("get saved skins: %w", err)
		}
	}
	defer rows.Close()

	items := make([]appskins.SavedSkin, 0, params.Limit)
	for rows.Next() {
		var item appskins.SavedSkin
		if err := rows.Scan(&item.MarketHashName, &item.DisplayName, &item.IconURL, &item.PageURL); err != nil {
			return appskins.SavedSkinsList{}, fmt.Errorf("scan saved skin: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return appskins.SavedSkinsList{}, fmt.Errorf("iterate saved skins: %w", err)
	}

	return appskins.SavedSkinsList{Items: items, TotalCount: totalCount, Offset: params.Offset, Limit: params.Limit}, nil
}
