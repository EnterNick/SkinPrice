package settings

type Storage interface {
	GetAppSettings() (AppSettings, error)
	SaveAppSettings(settings AppSettings) error
}
