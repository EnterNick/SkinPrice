package skins

import "SkinPrice/skinprice/internal/application"

type GetSavedSkins struct {
	SavedSkinsReader SavedSkinsReader
}

func (uc GetSavedSkins) Execute(params application.Pagination) (SavedSkinsList, error) {
	return uc.SavedSkinsReader.GetSavedList(&params)
}
