package skins

type SearchStorageSelector struct {
	DefaultStorage  NewSkinsStorage
	LisSkinsStorage NewSkinsStorage
}

func (s SearchStorageSelector) Get(source string) NewSkinsStorage {
	switch source {
	case "lisskins":
		if s.LisSkinsStorage != nil {
			return s.LisSkinsStorage
		}
	}
	return s.DefaultStorage
}
