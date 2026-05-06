package skins

import (
	"SkinPrice/skinprice/internal/application"
	"context"
)

type SearchNewSkins struct {
	NewSkinsStorage NewSkinsStorage
}

func (uc SearchNewSkins) Execute(ctx context.Context, criteria SearchCriteria, params application.Pagination) error {

}
