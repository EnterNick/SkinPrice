package skins

import "SkinPrice/skinprice/internal/application"

type NewSkinsStorage interface {
	GetList(criteria SearchCriteria, params *application.Pagination) (NewSkinsList, error)
}

type SkinSaver interface {
	Save(params SaveSkinParams) error
}
