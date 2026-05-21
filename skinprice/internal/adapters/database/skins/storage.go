package skins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/adapters/database/ent"
	entskin "SkinPrice/skinprice/internal/adapters/database/ent/skin"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"strings"
	"time"
)

type Storage struct {
	Conn *database.Connection
}

func (s *Storage) Save(ctx context.Context, params appskins.SaveSkinParams) (appskins.SaveSkinResult, error) {
	_, err := s.Conn.Client().Skin.Create().
		SetMarketHashName(params.MarketHashName).
		SetDisplayName(params.DisplayName).
		SetNameColor(params.NameColor).
		SetIconURL(params.IconURL).
		SetPageURL(params.PageURL).
		SetSteamPageURL(appskins.FirstNonEmpty(params.SteamPageURL, params.PageURL)).
		SetLisskinsPageURL(params.LisSkinsPageURL).
		SetCstmPageURL(params.CSTMPageURL).
		SetCurrency(appskins.NormalizeCurrencyCode("")).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) || isUniqueViolation(err) {
			return appskins.SaveSkinResult{Created: false}, nil
		}
		return appskins.SaveSkinResult{}, errx.E("skins.repository.save", errx.CodeInternal, "failed to save skin", err)
	}
	return appskins.SaveSkinResult{Created: true}, nil
}

func (s *Storage) GetSavedList(ctx context.Context, params *application.Pagination) (appskins.SavedSkinsList, error) {
	totalCount, err := s.Conn.Client().Skin.Query().Count(ctx)
	if err != nil {
		return appskins.SavedSkinsList{}, errx.E("skins.repository.list.count", errx.CodeInternal, "failed to count saved skins", err)
	}

	items, err := s.Conn.Client().Skin.Query().
		Order(ent.Desc(entskin.FieldID)).
		Limit(params.Limit).
		Offset(params.Offset).
		All(ctx)
	if err != nil {
		return appskins.SavedSkinsList{}, errx.E("skins.repository.list.query", errx.CodeInternal, "failed to load saved skins", err)
	}

	result := make([]appskins.SavedSkin, 0, len(items))
	for _, item := range items {
		result = append(result, mapSavedSkin(item))
	}
	return appskins.SavedSkinsList{Items: result, TotalCount: totalCount, Offset: params.Offset, Limit: params.Limit}, nil
}

func (s *Storage) GetSavedSkin(ctx context.Context, marketHashName string) (appskins.SavedSkin, error) {
	item, err := s.Conn.Client().Skin.Query().
		Where(entskin.MarketHashName(marketHashName)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return appskins.SavedSkin{}, errx.E("skins.repository.get.not_found", errx.CodeNotFound, "saved skin not found", err)
		}
		return appskins.SavedSkin{}, errx.E("skins.repository.get", errx.CodeInternal, "failed to load saved skin", err)
	}
	return mapSavedSkin(item), nil
}

func (s *Storage) ListSavedSkinNames(ctx context.Context) ([]string, error) {
	names, err := s.Conn.Client().Skin.Query().
		Order(ent.Desc(entskin.FieldID)).
		Select(entskin.FieldMarketHashName).
		Strings(ctx)
	if err != nil {
		return nil, errx.E("skins.repository.names", errx.CodeInternal, "failed to load saved skin names", err)
	}
	return names, nil
}

func (s *Storage) UpdateSavedSkinPrices(ctx context.Context, params appskins.UpdateSavedSkinPriceResult) error {
	now := latestUpdatedAt(params)
	count, err := s.Conn.Client().Skin.Update().
		Where(entskin.MarketHashName(params.MarketHashName)).
		SetPageURL(params.SteamPageURL).
		SetPriceText(params.SteamPriceText).
		SetSteamPageURL(params.SteamPageURL).
		SetSteamPriceText(params.SteamPriceText).
		SetNillableSteamUpdatedAt(nullableTime(params.SteamUpdatedAt)).
		SetLisskinsPageURL(params.LisSkinsPageURL).
		SetLisskinsPriceText(params.LisSkinsPriceText).
		SetNillableLisskinsUpdatedAt(nullableTime(params.LisSkinsUpdatedAt)).
		SetCstmPageURL(params.CSTMPageURL).
		SetCstmPriceText(params.CSTMPriceText).
		SetNillableCstmUpdatedAt(nullableTime(params.CSTMUpdatedAt)).
		SetCurrency(appskins.NormalizeCurrencyCode(params.Currency)).
		SetUpdatedAt(now).
		Save(ctx)
	if err != nil {
		return errx.E("skins.repository.update_prices", errx.CodeInternal, "failed to update saved skin price", err)
	}
	if count == 0 {
		return errx.E("skins.repository.update_prices.not_found", errx.CodeNotFound, "saved skin not found", nil)
	}
	return nil
}

func (s *Storage) DeleteSavedSkin(ctx context.Context, params appskins.DeleteSavedSkinParams) error {
	count, err := s.Conn.Client().Skin.Delete().
		Where(entskin.MarketHashName(params.MarketHashName)).
		Exec(ctx)
	if err != nil {
		return errx.E("skins.repository.delete", errx.CodeInternal, "failed to delete saved skin", err)
	}
	if count == 0 {
		return errx.E("skins.repository.delete.not_found", errx.CodeNotFound, "saved skin not found", nil)
	}
	return nil
}

func mapSavedSkin(item *ent.Skin) appskins.SavedSkin {
	saved := appskins.SavedSkin{
		MarketHashName:    item.MarketHashName,
		DisplayName:       item.DisplayName,
		NameColor:         item.NameColor,
		IconURL:           item.IconURL,
		SteamPageURL:      appskins.FirstNonEmpty(item.SteamPageURL, item.PageURL),
		SteamPriceText:    appskins.FirstNonEmpty(item.SteamPriceText, item.PriceText),
		LisSkinsPageURL:   item.LisskinsPageURL,
		LisSkinsPriceText: item.LisskinsPriceText,
		CSTMPageURL:       item.CstmPageURL,
		CSTMPriceText:     item.CstmPriceText,
		Currency:          appskins.NormalizeCurrencyCode(item.Currency),
	}
	if item.SteamUpdatedAt != nil {
		saved.SteamUpdatedAt = *item.SteamUpdatedAt
	} else if item.UpdatedAt != nil {
		saved.SteamUpdatedAt = *item.UpdatedAt
	}
	if item.LisskinsUpdatedAt != nil {
		saved.LisSkinsUpdatedAt = *item.LisskinsUpdatedAt
	}
	if item.CstmUpdatedAt != nil {
		saved.CSTMUpdatedAt = *item.CstmUpdatedAt
	}
	return saved
}

func nullableTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}

func latestUpdatedAt(params appskins.UpdateSavedSkinPriceResult) time.Time {
	latest := time.Time{}
	for _, value := range []time.Time{params.SteamUpdatedAt, params.LisSkinsUpdatedAt, params.CSTMUpdatedAt} {
		if !value.IsZero() && value.After(latest) {
			latest = value
		}
	}
	if latest.IsZero() {
		return time.Now().UTC()
	}
	return latest
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unique constraint") || strings.Contains(message, "duplicate key") || strings.Contains(message, "unique failed")
}
