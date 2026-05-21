package skins

import "context"

type SaveSkin struct {
	Repository    SavedSkinRepository
	LisSkinsPages MarketPageURLBuilder
	CSTMPages     MarketPageURLBuilder
}

func (uc SaveSkin) Execute(ctx context.Context, params SaveSkinParams) (SaveSkinResult, error) {
	params.SteamPageURL = FirstNonEmpty(params.SteamPageURL, params.PageURL)
	if uc.LisSkinsPages != nil {
		params.LisSkinsPageURL = FirstNonEmpty(params.LisSkinsPageURL, uc.LisSkinsPages.BuildMarketPageURL(params.MarketHashName))
	}
	if uc.CSTMPages != nil {
		params.CSTMPageURL = FirstNonEmpty(params.CSTMPageURL, uc.CSTMPages.BuildMarketPageURL(params.MarketHashName))
	}
	return uc.Repository.Save(ctx, params)
}
