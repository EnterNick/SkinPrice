package skins

import (
	"SkinPrice/skinprice/internal/application"
	"context"
)

type SearchNewSkins struct {
	Storage NewSkinsStorage
}

func (uc SearchNewSkins) Execute(ctx context.Context, criteria SearchCriteria, params application.Pagination) (NewSkinsList, error) {
	return uc.Storage.GetList(ctx, criteria, &params)
}
