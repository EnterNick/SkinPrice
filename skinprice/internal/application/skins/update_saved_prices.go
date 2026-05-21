package skins

import (
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"time"
)

type UpdateSavedSkinPrice struct {
	Repository SavedSkinRepository
	Collector  SavedSkinPriceCollector
}

func (uc UpdateSavedSkinPrice) Execute(ctx context.Context, params UpdateSavedSkinPriceParams) (UpdateSavedSkinPriceResult, error) {
	saved, err := uc.Repository.GetSavedSkin(ctx, params.MarketHashName)
	if err != nil {
		return UpdateSavedSkinPriceResult{}, err
	}

	if uc.Collector == nil {
		return UpdateSavedSkinPriceResult{}, errx.E("skins.update_one.no_collector", errx.CodeInternal, "no price collector configured", nil)
	}
	result, err := uc.Collector.Collect(ctx, saved, params)
	if err != nil {
		return UpdateSavedSkinPriceResult{}, err
	}

	if err := uc.Repository.UpdateSavedSkinPrices(ctx, result); err != nil {
		return UpdateSavedSkinPriceResult{}, err
	}
	return result, nil
}

type UpdateAllSavedSkinsPrices struct {
	Repository       SavedSkinRepository
	UpdateOne        SavedSkinPriceUpdater
	BatchUpdateDelay time.Duration
}

func (uc UpdateAllSavedSkinsPrices) Execute(ctx context.Context, params UpdateAllSavedSkinsPricesParams) (UpdateAllSavedSkinsPricesResult, error) {
	if uc.UpdateOne == nil {
		return UpdateAllSavedSkinsPricesResult{}, errx.E("skins.update_all.no_updater", errx.CodeInternal, "no price updater configured", nil)
	}

	names, err := uc.Repository.ListSavedSkinNames(ctx)
	if err != nil {
		return UpdateAllSavedSkinsPricesResult{}, err
	}

	failures := make([]UpdateSavedSkinPriceFailure, 0)
	updatedCount := 0
	currency := NormalizeCurrencyCode(params.Currency)
	for i, marketHashName := range names {
		if i > 0 && uc.BatchUpdateDelay > 0 {
			select {
			case <-ctx.Done():
				return UpdateAllSavedSkinsPricesResult{}, ctx.Err()
			case <-time.After(uc.BatchUpdateDelay):
			}
		}
		if _, err := uc.UpdateOne.Execute(ctx, UpdateSavedSkinPriceParams{
			MarketHashName: marketHashName,
			Currency:       currency,
		}); err != nil {
			failures = append(failures, UpdateSavedSkinPriceFailure{
				MarketHashName: marketHashName,
				Message:        err.Error(),
			})
			continue
		}
		updatedCount++
	}

	return UpdateAllSavedSkinsPricesResult{
		UpdatedCount: updatedCount,
		FailedCount:  len(failures),
		Failures:     failures,
	}, nil
}
