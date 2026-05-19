package skins

type SaveSkin struct {
	SkinSaver SkinSaver
}

func (uc SaveSkin) Execute(params SaveSkinParams) (SaveSkinResult, error) {
	return uc.SkinSaver.Save(params)
}
