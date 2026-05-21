package settings

import "context"

const (
	DefaultCurrency                   = "1"
	DefaultAutoRefreshEnabled         = true
	DefaultAutoRefreshIntervalSeconds = 30
	MinAutoRefreshIntervalSeconds     = 5
	DefaultSavedSkinsViewMode         = "table"
	DefaultFontFamily                 = "nunito"
	DefaultFontSizePx                 = 14
	MinFontSizePx                     = 10
	MaxFontSizePx                     = 28
)

type GetAppSettings struct {
	Storage Storage
}

func (uc GetAppSettings) Execute(ctx context.Context) (AppSettings, error) {
	settings, err := uc.Storage.GetAppSettings(ctx)
	if err != nil {
		return AppSettings{}, err
	}
	settings.Currency = normalizeCurrency(settings.Currency)
	settings.AutoRefreshEnabled = normalizeAutoRefreshEnabled(settings.AutoRefreshEnabled)
	settings.AutoRefreshIntervalSeconds = normalizeAutoRefreshIntervalSeconds(settings.AutoRefreshIntervalSeconds)
	settings.SavedSkinsViewMode = normalizeSavedSkinsViewMode(settings.SavedSkinsViewMode)
	settings.FontFamily = normalizeFontFamily(settings.FontFamily)
	settings.FontSizePx = normalizeFontSizePx(settings.FontSizePx)
	return settings, nil
}

type SaveAppSettings struct {
	Storage Storage
}

func (uc SaveAppSettings) Execute(ctx context.Context, settings AppSettings) error {
	settings.Currency = normalizeCurrency(settings.Currency)
	settings.AutoRefreshEnabled = normalizeAutoRefreshEnabled(settings.AutoRefreshEnabled)
	settings.AutoRefreshIntervalSeconds = normalizeAutoRefreshIntervalSeconds(settings.AutoRefreshIntervalSeconds)
	settings.SavedSkinsViewMode = normalizeSavedSkinsViewMode(settings.SavedSkinsViewMode)
	settings.FontFamily = normalizeFontFamily(settings.FontFamily)
	settings.FontSizePx = normalizeFontSizePx(settings.FontSizePx)
	return uc.Storage.SaveAppSettings(ctx, settings)
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

func normalizeAutoRefreshEnabled(value bool) bool {
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

func normalizeFontFamily(value string) string {
	switch value {
	case "inter", "system", "nunito", "roboto", "ibm-plex-sans", "manrope", "monocraft":
		return value
	default:
		return DefaultFontFamily
	}
}

func normalizeFontSizePx(value int) int {
	if value < MinFontSizePx || value > MaxFontSizePx {
		return DefaultFontSizePx
	}
	return value
}
