package skins

import (
	app "SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	"context"
)

type SearchNewSkinsUseCase interface {
	Execute(ctx context.Context, criteria appskins.SearchCriteria, params app.Pagination) (appskins.NewSkinsList, error)
}

type SaveSkinUseCase interface {
	Execute(ctx context.Context, params appskins.SaveSkinParams) (appskins.SaveSkinResult, error)
}

type GetSavedSkinsUseCase interface {
	Execute(ctx context.Context, params app.Pagination) (appskins.SavedSkinsList, error)
}
type UpdateSavedSkinPriceUseCase interface {
	Execute(ctx context.Context, params appskins.UpdateSavedSkinPriceParams) (appskins.UpdateSavedSkinPriceResult, error)
}
type UpdateAllSavedSkinsPricesUseCase interface {
	Execute(ctx context.Context, params appskins.UpdateAllSavedSkinsPricesParams) (appskins.UpdateAllSavedSkinsPricesResult, error)
}

type DeleteSavedSkinUseCase interface {
	Execute(ctx context.Context, params appskins.DeleteSavedSkinParams) error
}

type SaveLisSkinsTokenUseCase interface {
	Execute(ctx context.Context, token string) error
}

type HasLisSkinsTokenUseCase interface {
	Execute(ctx context.Context) (bool, error)
}

type ClearLisSkinsTokenUseCase interface {
	Execute(ctx context.Context) error
}

type Endpoints struct {
	searchNewSkinsUC            SearchNewSkinsUseCase
	saveSkinUC                  SaveSkinUseCase
	getSavedSkinsUC             GetSavedSkinsUseCase
	updateSavedSkinPriceUC      UpdateSavedSkinPriceUseCase
	updateAllSavedSkinsPricesUC UpdateAllSavedSkinsPricesUseCase
	deleteSavedSkinUC           DeleteSavedSkinUseCase
	saveLisSkinsTokenUC         SaveLisSkinsTokenUseCase
	hasLisSkinsTokenUC          HasLisSkinsTokenUseCase
	clearLisSkinsTokenUC        ClearLisSkinsTokenUseCase
}

type EndpointDeps struct {
	SearchNewSkins            SearchNewSkinsUseCase
	SaveSkin                  SaveSkinUseCase
	GetSavedSkins             GetSavedSkinsUseCase
	UpdateSavedSkinPrice      UpdateSavedSkinPriceUseCase
	UpdateAllSavedSkinsPrices UpdateAllSavedSkinsPricesUseCase
	DeleteSavedSkin           DeleteSavedSkinUseCase
	SaveLisSkinsToken         SaveLisSkinsTokenUseCase
	HasLisSkinsToken          HasLisSkinsTokenUseCase
	ClearLisSkinsToken        ClearLisSkinsTokenUseCase
}

func NewEndpoints(deps EndpointDeps) *Endpoints {
	return &Endpoints{
		searchNewSkinsUC:            deps.SearchNewSkins,
		saveSkinUC:                  deps.SaveSkin,
		getSavedSkinsUC:             deps.GetSavedSkins,
		updateSavedSkinPriceUC:      deps.UpdateSavedSkinPrice,
		updateAllSavedSkinsPricesUC: deps.UpdateAllSavedSkinsPrices,
		deleteSavedSkinUC:           deps.DeleteSavedSkin,
		saveLisSkinsTokenUC:         deps.SaveLisSkinsToken,
		hasLisSkinsTokenUC:          deps.HasLisSkinsToken,
		clearLisSkinsTokenUC:        deps.ClearLisSkinsToken,
	}
}

func (e *Endpoints) SearchNewSkins(ctx context.Context, filter SearchNewSkinsFilter) (NewSkinsResponse, error) {
	result, err := e.searchNewSkinsUC.Execute(ctx, appskins.SearchCriteria{
		MarketHashName:     filter.MarketHashName,
		SortColumn:         filter.SortColumn,
		SortDir:            filter.SortDir,
		PriceMin:           filter.PriceMin,
		PriceMax:           filter.PriceMax,
		SearchDescriptions: filter.SearchDescriptions,
		Types:              filter.Type,
		Weapons:            filter.Weapon,
		Rarities:           filter.Rarity,
		Exteriors:          filter.Exterior,
		ItemSets:           filter.ItemSet,
		ProPlayers:         filter.ProPlayer,
		StickerCapsules:    filter.StickerCapsule,
		TournamentTeams:    filter.TournamentTeam,
	}, app.Pagination{Limit: filter.Limit, Offset: filter.Offset, Cursor: filter.Cursor})
	if err != nil {
		return NewSkinsResponse{}, err
	}

	items := make([]NewSkinItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, NewSkinItem{
			MarketHashName: item.MarketHashName,
			DisplayName:    item.DisplayName,
			NameColor:      item.NameColor,
			SellListings:   item.SellListings,
			PriceCents:     item.PriceCents,
			PriceText:      item.PriceText,
			IconURL:        item.IconURL,
			PageURL:        item.PageURL,
		})
	}

	return NewSkinsResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Limit:      result.Limit,
		Offset:     result.Offset,
		NextCursor: result.NextCursor,
	}, nil
}

func (e *Endpoints) SaveSkin(ctx context.Context, payload SaveSkinRequest) (SaveSkinResponse, error) {
	result, err := e.saveSkinUC.Execute(ctx, appskins.SaveSkinParams{
		MarketHashName: payload.MarketHashName,
		DisplayName:    payload.DisplayName,
		NameColor:      payload.NameColor,
		IconURL:        payload.IconURL,
		PageURL:        payload.PageURL,
	})
	if err != nil {
		return SaveSkinResponse{}, err
	}
	return SaveSkinResponse{Created: result.Created}, nil
}

func (e *Endpoints) GetSavedSkins(ctx context.Context, filter GetSavedSkinsFilter) (SavedSkinsResponse, error) {
	result, err := e.getSavedSkinsUC.Execute(ctx, app.Pagination{Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return SavedSkinsResponse{}, err
	}

	items := make([]SavedSkinItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, SavedSkinItem{
			MarketHashName:    item.MarketHashName,
			DisplayName:       item.DisplayName,
			NameColor:         item.NameColor,
			IconURL:           item.IconURL,
			SteamPageURL:      item.SteamPageURL,
			SteamPriceText:    item.SteamPriceText,
			SteamUpdatedAt:    item.SteamUpdatedAt,
			LisSkinsPageURL:   item.LisSkinsPageURL,
			LisSkinsPriceText: item.LisSkinsPriceText,
			LisSkinsUpdatedAt: item.LisSkinsUpdatedAt,
			CSTMPageURL:       item.CSTMPageURL,
			CSTMPriceText:     item.CSTMPriceText,
			CSTMUpdatedAt:     item.CSTMUpdatedAt,
			Currency:          item.Currency,
		})
	}

	return SavedSkinsResponse{Items: items, TotalCount: result.TotalCount, Limit: result.Limit, Offset: result.Offset}, nil
}

func (e *Endpoints) UpdateSavedSkinPrice(ctx context.Context, payload UpdateSavedSkinPriceRequest) (UpdateSavedSkinPriceResponse, error) {
	result, err := e.updateSavedSkinPriceUC.Execute(ctx, appskins.UpdateSavedSkinPriceParams{MarketHashName: payload.MarketHashName, Currency: payload.Currency})
	if err != nil {
		return UpdateSavedSkinPriceResponse{}, err
	}
	return UpdateSavedSkinPriceResponse{
		MarketHashName:    result.MarketHashName,
		SteamPageURL:      result.SteamPageURL,
		SteamPriceText:    result.SteamPriceText,
		SteamUpdatedAt:    result.SteamUpdatedAt,
		LisSkinsPageURL:   result.LisSkinsPageURL,
		LisSkinsPriceText: result.LisSkinsPriceText,
		LisSkinsUpdatedAt: result.LisSkinsUpdatedAt,
		CSTMPageURL:       result.CSTMPageURL,
		CSTMPriceText:     result.CSTMPriceText,
		CSTMUpdatedAt:     result.CSTMUpdatedAt,
		Currency:          result.Currency,
	}, nil
}

func (e *Endpoints) UpdateAllSavedSkinsPrices(ctx context.Context, payload UpdateAllSavedSkinsPricesRequest) (UpdateAllSavedSkinsPricesResponse, error) {
	result, err := e.updateAllSavedSkinsPricesUC.Execute(ctx, appskins.UpdateAllSavedSkinsPricesParams{Currency: payload.Currency})
	if err != nil {
		return UpdateAllSavedSkinsPricesResponse{}, err
	}

	failures := make([]UpdateSavedSkinPriceFailure, 0, len(result.Failures))
	for _, failure := range result.Failures {
		failures = append(failures, UpdateSavedSkinPriceFailure{
			MarketHashName: failure.MarketHashName,
			Message:        failure.Message,
		})
	}

	return UpdateAllSavedSkinsPricesResponse{
		UpdatedCount: result.UpdatedCount,
		FailedCount:  result.FailedCount,
		Failures:     failures,
	}, nil
}

func (e *Endpoints) DeleteSavedSkin(ctx context.Context, payload DeleteSavedSkinRequest) error {
	return e.deleteSavedSkinUC.Execute(ctx, appskins.DeleteSavedSkinParams{MarketHashName: payload.MarketHashName})
}

func (e *Endpoints) SetLisSkinsToken(ctx context.Context, payload SetLisSkinsTokenRequest) error {
	return e.saveLisSkinsTokenUC.Execute(ctx, payload.Token)
}

func (e *Endpoints) GetLisSkinsTokenStatus(ctx context.Context) (LisSkinsTokenStatusResponse, error) {
	hasToken, err := e.hasLisSkinsTokenUC.Execute(ctx)
	if err != nil {
		return LisSkinsTokenStatusResponse{}, err
	}
	return LisSkinsTokenStatusResponse{HasToken: hasToken}, nil
}

func (e *Endpoints) ClearLisSkinsToken(ctx context.Context) error {
	return e.clearLisSkinsTokenUC.Execute(ctx)
}
