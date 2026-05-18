package skins

import "SkinPrice/skinprice/internal/application"

type SearchNewSkins struct {
	Storages NewSkinsStorageSelector
}

func (uc SearchNewSkins) Execute(criteria SearchCriteria, params application.Pagination) (NewSkinsList, error) {
	storage := uc.Storages.Get(criteria.Source)
	return storage.GetList(criteria, &params)
}
