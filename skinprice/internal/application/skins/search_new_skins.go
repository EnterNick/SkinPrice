package skins

import "SkinPrice/skinprice/internal/application"

type SearchNewSkins struct {
	NewSkinsStorage NewSkinsStorage
}

func (uc SearchNewSkins) Execute(criteria SearchCriteria, params application.Pagination) (NewSkinsList, error) {
	return uc.NewSkinsStorage.GetList(criteria, &params)
}
