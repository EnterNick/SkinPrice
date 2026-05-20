package settings

const (
	DefaultCurrency                   = "1"
	DefaultAutoRefreshIntervalSeconds = 30
	MinAutoRefreshIntervalSeconds     = 5
	DefaultSavedSkinsViewMode         = "table"
)

type GetAppSettings struct {
	Storage Storage
}

func (uc GetAppSettings) Execute() (AppSettings, error) {
	settings, err := uc.Storage.GetAppSettings()
	if err != nil {
		return AppSettings{}, err
	}
	settings.Currency = normalizeCurrency(settings.Currency)
	settings.AutoRefreshIntervalSeconds = normalizeAutoRefreshIntervalSeconds(settings.AutoRefreshIntervalSeconds)
	settings.SavedSkinsViewMode = normalizeSavedSkinsViewMode(settings.SavedSkinsViewMode)
	return settings, nil
}

type SaveAppSettings struct {
	Storage Storage
}

func (uc SaveAppSettings) Execute(settings AppSettings) error {
	settings.Currency = normalizeCurrency(settings.Currency)
	settings.AutoRefreshIntervalSeconds = normalizeAutoRefreshIntervalSeconds(settings.AutoRefreshIntervalSeconds)
	settings.SavedSkinsViewMode = normalizeSavedSkinsViewMode(settings.SavedSkinsViewMode)
	return uc.Storage.SaveAppSettings(settings)
}

func normalizeCurrency(value string) string {
	switch value {
	case "1", "USD":
		return "1"
	case "3", "EUR":
		return "3"
	case "5", "RUB":
		return "5"
	default:
		return DefaultCurrency
	}
}

func normalizeAutoRefreshIntervalSeconds(value int) int {
	if value < MinAutoRefreshIntervalSeconds {
		return DefaultAutoRefreshIntervalSeconds
	}
	return value
}

func normalizeSavedSkinsViewMode(value string) string {
	switch value {
	case "table", "cards":
		return value
	default:
		return DefaultSavedSkinsViewMode
	}
}
