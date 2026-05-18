package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	adaptersteam "SkinPrice/skinprice/internal/adapters/http/steam"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Storage struct {
	Conn         *database.Connection
	SteamStorage *adaptersteam.Storage
}

func ensureSkinsColumns(db *sql.DB) {
	ctx := context.Background()
	_, _ = db.ExecContext(ctx, `ALTER TABLE skins ADD COLUMN price_text TEXT NOT NULL DEFAULT ''`)
	_, _ = db.ExecContext(ctx, `ALTER TABLE skins ADD COLUMN currency TEXT NOT NULL DEFAULT 'USD'`)
	_, err := db.ExecContext(ctx, `ALTER TABLE skins ADD COLUMN updated_at TIMESTAMP`)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
		_, _ = db.ExecContext(ctx, `ALTER TABLE skins ADD COLUMN updated_at TEXT`)
	}
	_, _ = db.ExecContext(ctx, `UPDATE skins SET updated_at = CURRENT_TIMESTAMP WHERE updated_at IS NULL OR CAST(updated_at AS TEXT) = ''`)
}

func parseUpdatedAt(value string) time.Time {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed
		}
	}

	return time.Time{}
}

func (s *Storage) Save(params appskins.SaveSkinParams) error {
	db := s.Conn.DB()
	ctx := context.Background()
	ensureSkinsColumns(db)

	_, err := db.ExecContext(
		ctx,
		`INSERT INTO skins (market_hash_name, display_name, icon_url, page_url) VALUES ($1, $2, $3, $4) ON CONFLICT (market_hash_name) DO NOTHING`,
		params.MarketHashName,
		params.DisplayName,
		params.IconURL,
		params.PageURL,
	)
	if err != nil {
		_, sqliteErr := db.ExecContext(
			ctx,
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

func (s *Storage) GetSavedList(params *application.Pagination) (_ appskins.SavedSkinsList, err error) {
	db := s.Conn.DB()
	ctx := context.Background()
	ensureSkinsColumns(db)

	var totalCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM skins`).Scan(&totalCount); err != nil {
		return appskins.SavedSkinsList{}, fmt.Errorf("count skins: %w", err)
	}

	queryWithUpdatedAtPg := `SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency, CAST(updated_at AS TEXT) FROM skins ORDER BY id DESC LIMIT $1 OFFSET $2`
	queryWithUpdatedAtSqlite := `SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency, CAST(updated_at AS TEXT) FROM skins ORDER BY id DESC LIMIT ? OFFSET ?`
	queryFallbackPg := `SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency, '' FROM skins ORDER BY id DESC LIMIT $1 OFFSET $2`
	queryFallbackSqlite := `SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency, '' FROM skins ORDER BY id DESC LIMIT ? OFFSET ?`

	rows, err := db.QueryContext(ctx, queryWithUpdatedAtPg, params.Limit, params.Offset)
	if err != nil {
		rows, err = db.QueryContext(ctx, queryWithUpdatedAtSqlite, params.Limit, params.Offset)
		if err != nil {
			rows, err = db.QueryContext(ctx, queryFallbackPg, params.Limit, params.Offset)
			if err != nil {
				rows, err = db.QueryContext(ctx, queryFallbackSqlite, params.Limit, params.Offset)
				if err != nil {
					return appskins.SavedSkinsList{}, fmt.Errorf("get saved skins: %w", err)
				}
			}
		}
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close saved skins rows: %w", closeErr))
		}
	}()

	items := make([]appskins.SavedSkin, 0, params.Limit)
	for rows.Next() {
		var item appskins.SavedSkin
		var updatedAt string
		if err := rows.Scan(&item.MarketHashName, &item.DisplayName, &item.IconURL, &item.PageURL, &item.PriceText, &item.Currency, &updatedAt); err != nil {
			return appskins.SavedSkinsList{}, fmt.Errorf("scan saved skin: %w", err)
		}
		item.UpdatedAt = parseUpdatedAt(updatedAt)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return appskins.SavedSkinsList{}, fmt.Errorf("iterate saved skins: %w", err)
	}

	return appskins.SavedSkinsList{Items: items, TotalCount: totalCount, Offset: params.Offset, Limit: params.Limit}, nil
}

func (s *Storage) UpdateSavedSkinPrice(params appskins.UpdateSavedSkinPriceParams) error {
	ctx := context.Background()
	ensureSkinsColumns(s.Conn.DB())
	price, err := s.SteamStorage.GetByMarketHashName(params.MarketHashName, params.Currency)
	if err != nil {
		return fmt.Errorf("fetch price: %w", err)
	}

	_, err = s.Conn.DB().ExecContext(ctx, `UPDATE skins SET price_text = $1, currency = $2, updated_at = CURRENT_TIMESTAMP WHERE market_hash_name = $3`, price.PriceText, params.Currency, params.MarketHashName)
	if err != nil {
		_, err = s.Conn.DB().ExecContext(ctx, `UPDATE skins SET price_text = ?, currency = ?, updated_at = CURRENT_TIMESTAMP WHERE market_hash_name = ?`, price.PriceText, params.Currency, params.MarketHashName)
	}
	if err != nil {
		return fmt.Errorf("update saved skin price: %w", err)
	}
	return nil
}

func (s *Storage) UpdateAllSavedSkinsPrices(params appskins.UpdateAllSavedSkinsPricesParams) (err error) {
	ctx := context.Background()
	ensureSkinsColumns(s.Conn.DB())
	rows, err := s.Conn.DB().QueryContext(ctx, `SELECT market_hash_name FROM skins`)
	if err != nil {
		return fmt.Errorf("get saved skins names: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close saved skin names rows: %w", closeErr))
		}
	}()

	marketHashNames := make([]string, 0)
	for rows.Next() {
		var marketHashName string
		if err := rows.Scan(&marketHashName); err != nil {
			return fmt.Errorf("scan market hash name: %w", err)
		}
		marketHashNames = append(marketHashNames, marketHashName)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return rowsErr
	}

	if closeErr := rows.Close(); closeErr != nil {
		return fmt.Errorf("close saved skin names rows: %w", closeErr)
	}

	for _, marketHashName := range marketHashNames {
		if err := s.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{MarketHashName: marketHashName, Currency: params.Currency}); err != nil {
			return err
		}
	}

	return nil
}
