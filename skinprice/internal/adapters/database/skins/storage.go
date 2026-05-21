package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	adaptercstm "SkinPrice/skinprice/internal/adapters/http/cstm"
	adapterlisskins "SkinPrice/skinprice/internal/adapters/http/lisskins"
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
	"sync"
	"time"
)

type Storage struct {
	Conn             *database.Connection
	SteamStorage     marketPriceReader
	LisSkinsStorage  marketPriceReader
	CSTMStorage      marketPriceReader
	CSTMBaseURL      string
	BatchUpdateDelay time.Duration
	Logger           *slog.Logger
}

type marketPriceReader interface {
	GetByMarketHashName(marketHashName, currency string) (*appskins.NewSkin, error)
}

type savedSkinState struct {
	SteamPageURL      string
	SteamPriceText    string
	SteamUpdatedAt    time.Time
	LisSkinsPageURL   string
	LisSkinsPriceText string
	LisSkinsUpdatedAt time.Time
	CSTMPageURL       string
	CSTMPriceText     string
	CSTMUpdatedAt     time.Time
	Currency          string
}

const lisSkinsCurrency = "1"

func (s *Storage) Save(params appskins.SaveSkinParams) (appskins.SaveSkinResult, error) {
	logger := logx.WithComponent(s.Logger, "skins_storage")
	db := s.Conn.DB()
	ctx := context.Background()

	steamPageURL := params.PageURL
	lisSkinsPageURL := adapterlisskins.BuildMarketPageURL(params.MarketHashName)
	cstmPageURL := adaptercstm.BuildMarketPageURL(s.cstmBaseURL(), params.MarketHashName)
	insertQuery := `INSERT INTO skins (
		market_hash_name, display_name, name_color, icon_url, page_url, price_text,
		steam_page_url, steam_price_text, lisskins_page_url, lisskins_price_text, cstm_page_url, cstm_price_text, currency
	) VALUES (?, ?, ?, ?, ?, '', ?, '', ?, '', ?, '', '1')`
	if s.Conn.Dialect() == "postgres" {
		insertQuery = `INSERT INTO skins (
			market_hash_name, display_name, name_color, icon_url, page_url, price_text,
			steam_page_url, steam_price_text, lisskins_page_url, lisskins_price_text, cstm_page_url, cstm_price_text, currency
		) VALUES ($1, $2, $3, $4, $5, '', $6, '', $7, '', $8, '', '1')`
	}

	_, err := db.ExecContext(ctx, insertQuery, params.MarketHashName, params.DisplayName, params.NameColor, params.IconURL, steamPageURL, steamPageURL, lisSkinsPageURL, cstmPageURL)
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

	query := `SELECT
		market_hash_name, display_name, name_color, icon_url,
		page_url, price_text, steam_page_url, steam_price_text, steam_updated_at,
		lisskins_page_url, lisskins_price_text, lisskins_updated_at,
		cstm_page_url, cstm_price_text, cstm_updated_at, currency, updated_at
	FROM skins ORDER BY id DESC LIMIT ? OFFSET ?`
	if s.Conn.Dialect() == "postgres" {
		query = `SELECT
			market_hash_name, display_name, name_color, icon_url,
			page_url, price_text, steam_page_url, steam_price_text, steam_updated_at,
			lisskins_page_url, lisskins_price_text, lisskins_updated_at,
			cstm_page_url, cstm_price_text, cstm_updated_at, currency, updated_at
		FROM skins ORDER BY id DESC LIMIT $1 OFFSET $2`
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
		var legacyPageURL string
		var legacyPriceText string
		var steamUpdatedAt sql.NullTime
		var lisSkinsUpdatedAt sql.NullTime
		var cstmUpdatedAt sql.NullTime
		var legacyUpdatedAt sql.NullTime
		if err := rows.Scan(
			&item.MarketHashName,
			&item.DisplayName,
			&item.NameColor,
			&item.IconURL,
			&legacyPageURL,
			&legacyPriceText,
			&item.SteamPageURL,
			&item.SteamPriceText,
			&steamUpdatedAt,
			&item.LisSkinsPageURL,
			&item.LisSkinsPriceText,
			&lisSkinsUpdatedAt,
			&item.CSTMPageURL,
			&item.CSTMPriceText,
			&cstmUpdatedAt,
			&item.Currency,
			&legacyUpdatedAt,
		); err != nil {
			logger.Error("failed to scan saved skin", logx.ErrAttrs(err)...)
			return appskins.SavedSkinsList{}, errx.E("skins.storage.list.scan", errx.CodeInternal, "failed to scan saved skin", err)
		}
		item.Currency = normalizeCurrencyCode(item.Currency)
		if item.SteamPageURL == "" {
			item.SteamPageURL = legacyPageURL
		}
		if item.SteamPriceText == "" {
			item.SteamPriceText = legacyPriceText
		}
		if steamUpdatedAt.Valid {
			item.SteamUpdatedAt = steamUpdatedAt.Time
		} else if legacyUpdatedAt.Valid {
			item.SteamUpdatedAt = legacyUpdatedAt.Time
		}
		if lisSkinsUpdatedAt.Valid {
			item.LisSkinsUpdatedAt = lisSkinsUpdatedAt.Time
		}
		if cstmUpdatedAt.Valid {
			item.CSTMUpdatedAt = cstmUpdatedAt.Time
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
	state, err := s.loadSavedSkinState(ctx, params.MarketHashName)
	if err != nil {
		return appskins.UpdateSavedSkinPriceResult{}, err
	}

	updatedAny := false
	steamErr := error(nil)
	lisSkinsErr := error(nil)
	cstmErr := error(nil)
	var mu sync.Mutex
	var wg sync.WaitGroup

	if s.SteamStorage != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			price, fetchErr := s.SteamStorage.GetByMarketHashName(params.MarketHashName, currency)
			if fetchErr != nil {
				classifiedErr := classifySourceError("skins.storage.update_one.steam", "steam", fetchErr)
				logger.Warn("failed to fetch steam price",
					append([]any{
						slog.String("market_hash_name", params.MarketHashName),
						slog.String("currency", currency),
					}, logx.ErrAttrs(classifiedErr)...)...,
				)
				mu.Lock()
				steamErr = classifiedErr
				mu.Unlock()
				return
			}

			updatedAt := time.Now().UTC()
			mu.Lock()
			updatedAny = true
			state.SteamPageURL = firstNonEmpty(price.PageURL, state.SteamPageURL)
			state.SteamPriceText = normalizePriceText(price, currency)
			state.SteamUpdatedAt = updatedAt
			mu.Unlock()
		}()
	}

	if s.LisSkinsStorage != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			price, fetchErr := s.LisSkinsStorage.GetByMarketHashName(params.MarketHashName, lisSkinsCurrency)
			if fetchErr != nil {
				classifiedErr := classifySourceError("skins.storage.update_one.lisskins", "lisskins", fetchErr)
				logger.Warn("failed to fetch lisskins price",
					append([]any{
						slog.String("market_hash_name", params.MarketHashName),
						slog.String("currency", lisSkinsCurrency),
					}, logx.ErrAttrs(classifiedErr)...)...,
				)
				mu.Lock()
				lisSkinsErr = classifiedErr
				mu.Unlock()
				return
			}

			updatedAt := time.Now().UTC()
			mu.Lock()
			updatedAny = true
			state.LisSkinsPageURL = firstNonEmpty(price.PageURL, state.LisSkinsPageURL, adapterlisskins.BuildMarketPageURL(params.MarketHashName))
			state.LisSkinsPriceText = normalizePriceText(price, lisSkinsCurrency)
			state.LisSkinsUpdatedAt = updatedAt
			mu.Unlock()
		}()
	}

	if s.CSTMStorage != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			price, fetchErr := s.CSTMStorage.GetByMarketHashName(params.MarketHashName, currency)
			if fetchErr != nil {
				classifiedErr := classifySourceError("skins.storage.update_one.cstm", "cstm", fetchErr)
				logger.Warn("failed to fetch cstm price",
					append([]any{
						slog.String("market_hash_name", params.MarketHashName),
						slog.String("currency", currency),
					}, logx.ErrAttrs(classifiedErr)...)...,
				)
				mu.Lock()
				cstmErr = classifiedErr
				mu.Unlock()
				return
			}

			updatedAt := time.Now().UTC()
			mu.Lock()
			updatedAny = true
			state.CSTMPageURL = firstNonEmpty(price.PageURL, state.CSTMPageURL, adaptercstm.BuildMarketPageURL(s.cstmBaseURL(), params.MarketHashName))
			state.CSTMPriceText = normalizePriceText(price, currency)
			state.CSTMUpdatedAt = updatedAt
			mu.Unlock()
		}()
	}

	wg.Wait()

	if !updatedAny {
		if steamErr != nil {
			return appskins.UpdateSavedSkinPriceResult{}, steamErr
		}
		if lisSkinsErr != nil {
			return appskins.UpdateSavedSkinPriceResult{}, lisSkinsErr
		}
		if cstmErr != nil {
			return appskins.UpdateSavedSkinPriceResult{}, cstmErr
		}
		return appskins.UpdateSavedSkinPriceResult{}, errx.E("skins.storage.update_one.no_sources", errx.CodeInternal, "no price sources configured", nil)
	}

	state.Currency = currency
	query := `UPDATE skins SET
		page_url = ?,
		price_text = ?,
		steam_page_url = ?,
		steam_price_text = ?,
		steam_updated_at = ?,
		lisskins_page_url = ?,
		lisskins_price_text = ?,
		lisskins_updated_at = ?,
		cstm_page_url = ?,
		cstm_price_text = ?,
		cstm_updated_at = ?,
		currency = ?,
		updated_at = ?
	WHERE market_hash_name = ?`
	if s.Conn.Dialect() == "postgres" {
		query = `UPDATE skins SET
			page_url = $1,
			price_text = $2,
			steam_page_url = $3,
			steam_price_text = $4,
			steam_updated_at = $5,
			lisskins_page_url = $6,
			lisskins_price_text = $7,
			lisskins_updated_at = $8,
			cstm_page_url = $9,
			cstm_price_text = $10,
			cstm_updated_at = $11,
			currency = $12,
			updated_at = $13
		WHERE market_hash_name = $14`
	}

	result, err := s.Conn.DB().ExecContext(
		ctx,
		query,
		state.SteamPageURL,
		state.SteamPriceText,
		state.SteamPageURL,
		state.SteamPriceText,
		nullTime(state.SteamUpdatedAt),
		state.LisSkinsPageURL,
		state.LisSkinsPriceText,
		nullTime(state.LisSkinsUpdatedAt),
		state.CSTMPageURL,
		state.CSTMPriceText,
		nullTime(state.CSTMUpdatedAt),
		state.Currency,
		nullTime(state.SteamUpdatedAt),
		params.MarketHashName,
	)
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
		slog.String("steam_price_text", state.SteamPriceText),
		slog.String("lisskins_price_text", state.LisSkinsPriceText),
		slog.String("cstm_price_text", state.CSTMPriceText),
	)
	return appskins.UpdateSavedSkinPriceResult{
		MarketHashName:    params.MarketHashName,
		SteamPageURL:      state.SteamPageURL,
		SteamPriceText:    state.SteamPriceText,
		SteamUpdatedAt:    state.SteamUpdatedAt,
		LisSkinsPageURL:   state.LisSkinsPageURL,
		LisSkinsPriceText: state.LisSkinsPriceText,
		LisSkinsUpdatedAt: state.LisSkinsUpdatedAt,
		CSTMPageURL:       state.CSTMPageURL,
		CSTMPriceText:     state.CSTMPriceText,
		CSTMUpdatedAt:     state.CSTMUpdatedAt,
		Currency:          state.Currency,
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

func (s *Storage) loadSavedSkinState(ctx context.Context, marketHashName string) (savedSkinState, error) {
	query := `SELECT
		page_url, price_text, steam_page_url, steam_price_text, steam_updated_at,
		lisskins_page_url, lisskins_price_text, lisskins_updated_at,
		cstm_page_url, cstm_price_text, cstm_updated_at, currency, updated_at
	FROM skins WHERE market_hash_name = ?`
	if s.Conn.Dialect() == "postgres" {
		query = `SELECT
			page_url, price_text, steam_page_url, steam_price_text, steam_updated_at,
			lisskins_page_url, lisskins_price_text, lisskins_updated_at,
			cstm_page_url, cstm_price_text, cstm_updated_at, currency, updated_at
		FROM skins WHERE market_hash_name = $1`
	}

	var state savedSkinState
	var legacyPageURL string
	var legacyPriceText string
	var steamUpdatedAt sql.NullTime
	var lisSkinsUpdatedAt sql.NullTime
	var cstmUpdatedAt sql.NullTime
	var legacyUpdatedAt sql.NullTime
	err := s.Conn.DB().QueryRowContext(ctx, query, marketHashName).Scan(
		&legacyPageURL,
		&legacyPriceText,
		&state.SteamPageURL,
		&state.SteamPriceText,
		&steamUpdatedAt,
		&state.LisSkinsPageURL,
		&state.LisSkinsPriceText,
		&lisSkinsUpdatedAt,
		&state.CSTMPageURL,
		&state.CSTMPriceText,
		&cstmUpdatedAt,
		&state.Currency,
		&legacyUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return savedSkinState{}, errx.E("skins.storage.update_one.not_found", errx.CodeNotFound, "saved skin not found", err)
		}
		return savedSkinState{}, errx.E("skins.storage.update_one.load", errx.CodeInternal, "failed to load saved skin state", err)
	}

	state.SteamPageURL = firstNonEmpty(state.SteamPageURL, legacyPageURL)
	state.SteamPriceText = firstNonEmpty(state.SteamPriceText, legacyPriceText)
	state.LisSkinsPageURL = firstNonEmpty(state.LisSkinsPageURL, adapterlisskins.BuildMarketPageURL(marketHashName))
	state.CSTMPageURL = firstNonEmpty(state.CSTMPageURL, adaptercstm.BuildMarketPageURL(s.cstmBaseURL(), marketHashName))
	state.Currency = normalizeCurrencyCode(state.Currency)
	if steamUpdatedAt.Valid {
		state.SteamUpdatedAt = steamUpdatedAt.Time
	} else if legacyUpdatedAt.Valid {
		state.SteamUpdatedAt = legacyUpdatedAt.Time
	}
	if lisSkinsUpdatedAt.Valid {
		state.LisSkinsUpdatedAt = lisSkinsUpdatedAt.Time
	}
	if cstmUpdatedAt.Valid {
		state.CSTMUpdatedAt = cstmUpdatedAt.Time
	}
	return state, nil
}

func nullTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value, Valid: true}
}

func normalizePriceText(price *appskins.NewSkin, currency string) string {
	if price == nil {
		return ""
	}
	if price.PriceCents != nil {
		return formatPriceText(*price.PriceCents, currency)
	}
	return price.PriceText
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (s *Storage) cstmBaseURL() string {
	if strings.TrimSpace(s.CSTMBaseURL) != "" {
		return s.CSTMBaseURL
	}
	return "https://market.csgo.com"
}

func classifySourceError(op, source string, err error) error {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "context deadline exceeded"):
		return errx.E(op, errx.CodeTimeout, source+" request timed out", err)
	case strings.Contains(message, "bad status"), strings.Contains(message, "unsuccessful"):
		return errx.E(op, errx.CodeUnavailable, source+" is temporarily unavailable", err)
	case strings.Contains(message, "not found"):
		return errx.E(op, errx.CodeNotFound, "skin not found on "+source, err)
	case strings.Contains(message, "decode"):
		return errx.E(op, errx.CodeExternal, "failed to decode "+source+" response", err)
	default:
		return errx.E(op, errx.CodeExternal, source+" request failed", err)
	}
}
