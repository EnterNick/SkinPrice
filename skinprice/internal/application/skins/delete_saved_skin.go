package skins

import "context"

type DeleteSavedSkin struct {
	Repository SavedSkinRepository
}

func (uc DeleteSavedSkin) Execute(ctx context.Context, params DeleteSavedSkinParams) error {
	return uc.Repository.DeleteSavedSkin(ctx, params)
}
