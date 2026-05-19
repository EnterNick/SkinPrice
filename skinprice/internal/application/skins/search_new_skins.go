package skins

import "SkinPrice/skinprice/internal/application"

type SearchNewSkins struct {
	Storage NewSkinsStorage
}

func (uc SearchNewSkins) Execute(criteria SearchCriteria, params application.Pagination) (NewSkinsList, error) {
	return uc.Storage.GetList(criteria, &params)
}
