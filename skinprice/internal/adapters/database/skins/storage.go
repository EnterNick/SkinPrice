package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	adaptersteam "SkinPrice/skinprice/internal/adapters/http/steam"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"fmt"
)

type Storage struct {
	Conn         *database.Connection
	SteamStorage *adaptersteam.Storage
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

func (s *Storage) GetSavedList(params *application.Pagination) (appskins.SavedSkinsList, error) {
	db := s.Conn.DB()
	_, _ = db.Exec(`ALTER TABLE skins ADD COLUMN price_text TEXT NOT NULL DEFAULT ''`)
	_, _ = db.Exec(`ALTER TABLE skins ADD COLUMN currency TEXT NOT NULL DEFAULT 'USD'`)

	var totalCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM skins`).Scan(&totalCount); err != nil {
		return appskins.SavedSkinsList{}, fmt.Errorf("count skins: %w", err)
	}

	rows, err := db.Query(
		`SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency FROM skins ORDER BY id DESC LIMIT $1 OFFSET $2`,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		rows, err = db.Query(
			`SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency FROM skins ORDER BY id DESC LIMIT ? OFFSET ?`,
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
		if err := rows.Scan(&item.MarketHashName, &item.DisplayName, &item.IconURL, &item.PageURL, &item.PriceText, &item.Currency); err != nil {
			return appskins.SavedSkinsList{}, fmt.Errorf("scan saved skin: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return appskins.SavedSkinsList{}, fmt.Errorf("iterate saved skins: %w", err)
	}

	return appskins.SavedSkinsList{Items: items, TotalCount: totalCount, Offset: params.Offset, Limit: params.Limit}, nil
}

func (s *Storage) UpdateSavedSkinPrice(params appskins.UpdateSavedSkinPriceParams) error {
	price, err := s.SteamStorage.GetByMarketHashName(params.MarketHashName, params.Currency)
	if err != nil {
		return fmt.Errorf("fetch price: %w", err)
	}

	_, err = s.Conn.DB().Exec(`UPDATE skins SET price_text = $1, currency = $2 WHERE market_hash_name = $3`, price.PriceText, params.Currency, params.MarketHashName)
	if err != nil {
		_, err = s.Conn.DB().Exec(`UPDATE skins SET price_text = ?, currency = ? WHERE market_hash_name = ?`, price.PriceText, params.Currency, params.MarketHashName)
	}
	if err != nil {
		return fmt.Errorf("update saved skin price: %w", err)
	}
	return nil
}

func (s *Storage) UpdateAllSavedSkinsPrices(params appskins.UpdateAllSavedSkinsPricesParams) error {
	rows, err := s.Conn.DB().Query(`SELECT market_hash_name FROM skins`)
	if err != nil {
		return fmt.Errorf("get saved skins names: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var marketHashName string
		if err := rows.Scan(&marketHashName); err != nil {
			return fmt.Errorf("scan market hash name: %w", err)
		}
		if err := s.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{MarketHashName: marketHashName, Currency: params.Currency}); err != nil {
			return err
		}
	}

	return rows.Err()
}
