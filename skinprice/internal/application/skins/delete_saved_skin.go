package skins

type DeleteSavedSkin struct {
	Deleter SavedSkinDeleter
}

func (uc DeleteSavedSkin) Execute(params DeleteSavedSkinParams) error {
	return uc.Deleter.DeleteSavedSkin(params)
}
