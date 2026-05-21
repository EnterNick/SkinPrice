package skins

import "SkinPrice/skinprice/internal/application"
import "context"

type GetSavedSkins struct {
	Repository SavedSkinRepository
}

func (uc GetSavedSkins) Execute(ctx context.Context, params application.Pagination) (SavedSkinsList, error) {
	return uc.Repository.GetSavedList(ctx, &params)
}
