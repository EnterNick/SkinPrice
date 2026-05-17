package skins

import (
	app "SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
)

type SearchNewSkinsUseCase interface {
	Execute(criteria appskins.SearchCriteria, params app.Pagination) (appskins.NewSkinsList, error)
}

type SaveSkinUseCase interface {
	Execute(params appskins.SaveSkinParams) error
}

type GetSavedSkinsUseCase interface {
	Execute(params app.Pagination) (appskins.SavedSkinsList, error)
}

type Endpoints struct {
	searchNewSkinsUC SearchNewSkinsUseCase
	saveSkinUC       SaveSkinUseCase
	getSavedSkinsUC  GetSavedSkinsUseCase
}

func NewEndpoints(searchNewSkinsUC SearchNewSkinsUseCase, saveSkinUC SaveSkinUseCase, getSavedSkinsUC GetSavedSkinsUseCase) *Endpoints {
	return &Endpoints{searchNewSkinsUC: searchNewSkinsUC, saveSkinUC: saveSkinUC, getSavedSkinsUC: getSavedSkinsUC}
}

func (e *Endpoints) SearchNewSkins(filter SearchNewSkinsFilter) (NewSkinsResponse, error) {
	result, err := e.searchNewSkinsUC.Execute(appskins.SearchCriteria{MarketHashName: filter.MarketHashName}, app.Pagination{Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return NewSkinsResponse{}, err
	}

	items := make([]NewSkinItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, NewSkinItem{
			MarketHashName: item.MarketHashName,
			DisplayName:    item.DisplayName,
			SellListings:   item.SellListings,
			PriceCents:     item.PriceCents,
			PriceText:      item.PriceText,
			IconURL:        item.IconURL,
			PageURL:        item.PageURL,
		})
	}

	return NewSkinsResponse{Items: items, TotalCount: result.TotalCount, Limit: result.Limit, Offset: result.Offset}, nil
}

func (e *Endpoints) SaveSkin(payload SaveSkinRequest) error {
	return e.saveSkinUC.Execute(appskins.SaveSkinParams{
		MarketHashName: payload.MarketHashName,
		DisplayName:    payload.DisplayName,
		IconURL:        payload.IconURL,
		PageURL:        payload.PageURL,
	})
}

func (e *Endpoints) GetSavedSkins(filter GetSavedSkinsFilter) (SavedSkinsResponse, error) {
	result, err := e.getSavedSkinsUC.Execute(app.Pagination{Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return SavedSkinsResponse{}, err
	}

	items := make([]SavedSkinItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, SavedSkinItem{
			MarketHashName: item.MarketHashName,
			DisplayName:    item.DisplayName,
			IconURL:        item.IconURL,
			PageURL:        item.PageURL,
		})
	}

	return SavedSkinsResponse{Items: items, TotalCount: result.TotalCount, Limit: result.Limit, Offset: result.Offset}, nil
}
