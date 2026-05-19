package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/shared/errx"
	"SkinPrice/skinprice/internal/shared/logx"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type Storage struct {
	Conn             *database.Connection
	SteamStorage     steamPriceReader
	BatchUpdateDelay time.Duration
	Logger           *slog.Logger
}

type steamPriceReader interface {
	GetByMarketHashName(marketHashName, currency string) (*appskins.NewSkin, error)
}

func (s *Storage) Save(params appskins.SaveSkinParams) (appskins.SaveSkinResult, error) {
	logger := logx.WithComponent(s.Logger, "skins_storage")
	db := s.Conn.DB()
	ctx := context.Background()

	insertQuery := `INSERT INTO skins (market_hash_name, display_name, icon_url, page_url, price_text, currency) VALUES (?, ?, ?, ?, '', '1')`
	if s.Conn.Dialect() == "postgres" {
		insertQuery = `INSERT INTO skins (market_hash_name, display_name, icon_url, page_url, price_text, currency) VALUES ($1, $2, $3, $4, '', '1')`
	}

	_, err := db.ExecContext(ctx, insertQuery, params.MarketHashName, params.DisplayName, params.IconURL, params.PageURL)
	if err != nil {
		if isUniqueViolation(err) {
			logger.Info("skin already exists", slog.String("market_hash_name", params.MarketHashName))
			return appskins.SaveSkinResult{Created: false}, nil
		}
		logger.Error("failed to save skin",
			append([]any{slog.String("market_hash_name", params.MarketHashName)}, logx.ErrAttrs(err)...)...,
		)
		return appskins.SaveSkinResult{}, errx.E("skins.storage.save", errx.CodeInternal, "failed to save skin", err)
	}

	logger.Info("skin saved", slog.String("market_hash_name", params.MarketHashName))
	return appskins.SaveSkinResult{Created: true}, nil
}

func (s *Storage) GetSavedList(params *application.Pagination) (_ appskins.SavedSkinsList, err error) {
	logger := logx.WithComponent(s.Logger, "skins_storage")
	db := s.Conn.DB()
	ctx := context.Background()

	var totalCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM skins`).Scan(&totalCount); err != nil {
		logger.Error("failed to count saved skins", logx.ErrAttrs(err)...)
		return appskins.SavedSkinsList{}, errx.E("skins.storage.list.count", errx.CodeInternal, "failed to count saved skins", err)
	}

	query := `SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency, updated_at FROM skins ORDER BY id DESC LIMIT ? OFFSET ?`
	if s.Conn.Dialect() == "postgres" {
		query = `SELECT market_hash_name, display_name, icon_url, page_url, price_text, currency, updated_at FROM skins ORDER BY id DESC LIMIT $1 OFFSET $2`
	}

	rows, err := db.QueryContext(ctx, query, params.Limit, params.Offset)
	if err != nil {
		logger.Error("failed to query saved skins", logx.ErrAttrs(err)...)
		return appskins.SavedSkinsList{}, errx.E("skins.storage.list.query", errx.CodeInternal, "failed to load saved skins", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close saved skins rows: %w", closeErr))
		}
	}()

	items := make([]appskins.SavedSkin, 0, params.Limit)
	for rows.Next() {
		var item appskins.SavedSkin
		var updatedAt sql.NullTime
		if err := rows.Scan(&item.MarketHashName, &item.DisplayName, &item.IconURL, &item.PageURL, &item.PriceText, &item.Currency, &updatedAt); err != nil {
			logger.Error("failed to scan saved skin", logx.ErrAttrs(err)...)
			return appskins.SavedSkinsList{}, errx.E("skins.storage.list.scan", errx.CodeInternal, "failed to scan saved skin", err)
		}
		item.Currency = normalizeCurrencyCode(item.Currency)
		if updatedAt.Valid {
			item.UpdatedAt = updatedAt.Time
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		logger.Error("failed while iterating saved skins", logx.ErrAttrs(err)...)
		return appskins.SavedSkinsList{}, errx.E("skins.storage.list.rows", errx.CodeInternal, "failed to iterate saved skins", err)
	}

	logger.Info("loaded saved skins",
		slog.Int("count", len(items)),
		slog.Int("total_count", totalCount),
		slog.Int("limit", params.Limit),
		slog.Int("offset", params.Offset),
	)
	return appskins.SavedSkinsList{Items: items, TotalCount: totalCount, Offset: params.Offset, Limit: params.Limit}, nil
}

func (s *Storage) UpdateSavedSkinPrice(params appskins.UpdateSavedSkinPriceParams) (appskins.UpdateSavedSkinPriceResult, error) {
	logger := logx.WithComponent(s.Logger, "skins_storage")
	ctx := context.Background()
	currency := normalizeCurrencyCode(params.Currency)
	price, err := s.SteamStorage.GetByMarketHashName(params.MarketHashName, currency)
	if err != nil {
		logger.Error("failed to fetch latest price",
			append([]any{
				slog.String("market_hash_name", params.MarketHashName),
				slog.String("currency", currency),
			}, logx.ErrAttrs(err)...)...,
		)
		return appskins.UpdateSavedSkinPriceResult{}, classifySteamError("skins.storage.update_one.fetch", err)
	}

	query := `UPDATE skins SET price_text = ?, currency = ?, updated_at = ? WHERE market_hash_name = ?`
	if s.Conn.Dialect() == "postgres" {
		query = `UPDATE skins SET price_text = $1, currency = $2, updated_at = $3 WHERE market_hash_name = $4`
	}

	updatedAt := time.Now().UTC()
	priceText := price.PriceText
	if price.PriceCents != nil {
		priceText = formatPriceText(*price.PriceCents, currency)
	}

	result, err := s.Conn.DB().ExecContext(ctx, query, priceText, currency, updatedAt, params.MarketHashName)
	if err != nil {
		logger.Error("failed to update saved skin price",
			append([]any{slog.String("market_hash_name", params.MarketHashName)}, logx.ErrAttrs(err)...)...,
		)
		return appskins.UpdateSavedSkinPriceResult{}, errx.E("skins.storage.update_one.exec", errx.CodeInternal, "failed to update saved skin price", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("failed to inspect update result",
			append([]any{slog.String("market_hash_name", params.MarketHashName)}, logx.ErrAttrs(err)...)...,
		)
		return appskins.UpdateSavedSkinPriceResult{}, errx.E("skins.storage.update_one.rows", errx.CodeInternal, "failed to inspect update result", err)
	}
	if rowsAffected == 0 {
		logger.Warn("saved skin not found during price update", slog.String("market_hash_name", params.MarketHashName))
		return appskins.UpdateSavedSkinPriceResult{}, errx.E("skins.storage.update_one.not_found", errx.CodeNotFound, "saved skin not found", nil)
	}

	logger.Info("saved skin price updated",
		slog.String("market_hash_name", params.MarketHashName),
		slog.String("currency", currency),
		slog.String("price_text", priceText),
	)
	return appskins.UpdateSavedSkinPriceResult{
		MarketHashName: params.MarketHashName,
		PriceText:      priceText,
		Currency:       currency,
		UpdatedAt:      updatedAt,
	}, nil
}

func (s *Storage) UpdateAllSavedSkinsPrices(params appskins.UpdateAllSavedSkinsPricesParams) (result appskins.UpdateAllSavedSkinsPricesResult, err error) {
	logger := logx.WithComponent(s.Logger, "skins_storage")
	ctx := context.Background()
	query := `SELECT market_hash_name FROM skins ORDER BY id DESC`
	rows, err := s.Conn.DB().QueryContext(ctx, query)
	if err != nil {
		logger.Error("failed to query saved skin names", logx.ErrAttrs(err)...)
		return appskins.UpdateAllSavedSkinsPricesResult{}, errx.E("skins.storage.update_all.query", errx.CodeInternal, "failed to load saved skin names", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close saved skin names rows: %w", closeErr))
		}
	}()

	names := make([]string, 0)
	for rows.Next() {
		var marketHashName string
		if err := rows.Scan(&marketHashName); err != nil {
			logger.Error("failed to scan saved skin name", logx.ErrAttrs(err)...)
			return appskins.UpdateAllSavedSkinsPricesResult{}, errx.E("skins.storage.update_all.scan", errx.CodeInternal, "failed to scan saved skin name", err)
		}
		names = append(names, marketHashName)
	}
	if err := rows.Err(); err != nil {
		logger.Error("failed while iterating saved skin names", logx.ErrAttrs(err)...)
		return appskins.UpdateAllSavedSkinsPricesResult{}, errx.E("skins.storage.update_all.rows", errx.CodeInternal, "failed to iterate saved skin names", err)
	}

	failures := make([]appskins.UpdateSavedSkinPriceFailure, 0)
	updatedCount := 0
	currency := normalizeCurrencyCode(params.Currency)
	for i, marketHashName := range names {
		if i > 0 && s.BatchUpdateDelay > 0 {
			time.Sleep(s.BatchUpdateDelay)
		}
		if _, err := s.UpdateSavedSkinPrice(appskins.UpdateSavedSkinPriceParams{
			MarketHashName: marketHashName,
			Currency:       currency,
		}); err != nil {
			logger.Warn("failed to update one saved skin in batch",
				append([]any{slog.String("market_hash_name", marketHashName)}, logx.ErrAttrs(err)...)...,
			)
			failures = append(failures, appskins.UpdateSavedSkinPriceFailure{
				MarketHashName: marketHashName,
				Message:        err.Error(),
			})
			continue
		}
		updatedCount++
	}

	logger.Info("bulk saved skin price update completed",
		slog.Int("requested_count", len(names)),
		slog.Int("updated_count", updatedCount),
		slog.Int("failed_count", len(failures)),
		slog.String("currency", currency),
		slog.Duration("delay_between_requests", s.BatchUpdateDelay),
	)
	return appskins.UpdateAllSavedSkinsPricesResult{
		UpdatedCount: updatedCount,
		FailedCount:  len(failures),
		Failures:     failures,
	}, nil
}

func (s *Storage) DeleteSavedSkin(params appskins.DeleteSavedSkinParams) error {
	logger := logx.WithComponent(s.Logger, "skins_storage")
	ctx := context.Background()
	query := `DELETE FROM skins WHERE market_hash_name = ?`
	if s.Conn.Dialect() == "postgres" {
		query = `DELETE FROM skins WHERE market_hash_name = $1`
	}

	result, err := s.Conn.DB().ExecContext(ctx, query, params.MarketHashName)
	if err != nil {
		logger.Error("failed to delete saved skin",
			append([]any{slog.String("market_hash_name", params.MarketHashName)}, logx.ErrAttrs(err)...)...,
		)
		return errx.E("skins.storage.delete", errx.CodeInternal, "failed to delete saved skin", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("failed to inspect delete result",
			append([]any{slog.String("market_hash_name", params.MarketHashName)}, logx.ErrAttrs(err)...)...,
		)
		return errx.E("skins.storage.delete.rows", errx.CodeInternal, "failed to inspect delete result", err)
	}
	if rowsAffected == 0 {
		logger.Warn("saved skin not found during delete", slog.String("market_hash_name", params.MarketHashName))
		return errx.E("skins.storage.delete.not_found", errx.CodeNotFound, "saved skin not found", nil)
	}
	logger.Info("saved skin deleted", slog.String("market_hash_name", params.MarketHashName))
	return nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unique constraint") || strings.Contains(message, "duplicate key") || strings.Contains(message, "unique failed")
}

func normalizeCurrencyCode(currency string) string {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "1", "USD":
		return "1"
	case "3", "EUR":
		return "3"
	case "5", "RUB":
		return "5"
	default:
		return "1"
	}
}

func formatPriceText(priceCents int64, currency string) string {
	sign := ""
	if priceCents < 0 {
		sign = "-"
		priceCents = -priceCents
	}

	whole := priceCents / 100
	fraction := priceCents % 100

	switch currency {
	case "3":
		return fmt.Sprintf("%s€%d.%02d", sign, whole, fraction)
	case "5":
		if fraction == 0 {
			return fmt.Sprintf("%s%d ₽", sign, whole)
		}
		return fmt.Sprintf("%s%d.%02d ₽", sign, whole, fraction)
	default:
		return fmt.Sprintf("%s$%d.%02d", sign, whole, fraction)
	}
}

func classifySteamError(op string, err error) error {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "context deadline exceeded"):
		return errx.E(op, errx.CodeTimeout, "steam request timed out", err)
	case strings.Contains(message, "bad status"), strings.Contains(message, "unsuccessful"):
		return errx.E(op, errx.CodeUnavailable, "steam is temporarily unavailable", err)
	case strings.Contains(message, "not found"):
		return errx.E(op, errx.CodeNotFound, "skin not found on steam", err)
	case strings.Contains(message, "decode"):
		return errx.E(op, errx.CodeExternal, "failed to decode steam response", err)
	default:
		return errx.E(op, errx.CodeExternal, "steam request failed", err)
	}
}
