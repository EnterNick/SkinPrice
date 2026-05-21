package main

import (
	adapterdbappsettings "SkinPrice/skinprice/internal/adapters/database/appsettings"
	adapterdbskins "SkinPrice/skinprice/internal/adapters/database/skins"
	adapterdbsourcestate "SkinPrice/skinprice/internal/adapters/database/sourcestate"
	adaptercstm "SkinPrice/skinprice/internal/adapters/http/cstm"
	adapterlisskins "SkinPrice/skinprice/internal/adapters/http/lisskins"
	adaptersteam "SkinPrice/skinprice/internal/adapters/http/steam"
	appsettings "SkinPrice/skinprice/internal/application/settings"
	"SkinPrice/skinprice/internal/application/skins"
	"SkinPrice/skinprice/internal/config"
	presentersettings "SkinPrice/skinprice/internal/presenters/settings"
	presenterskins "SkinPrice/skinprice/internal/presenters/skins"
	"SkinPrice/skinprice/internal/shared/errx"
	"SkinPrice/skinprice/internal/shared/logx"
	sharedcrypto "SkinPrice/skinprice/internal/shared/utils/crypto"
	"log/slog"
)

func (a *App) registerRoutes() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	steamStorage := &adaptersteam.Storage{Client: adaptersteam.NewSteamClient(cfg), BaseURL: cfg.SteamBaseURL, Logger: a.logger}
	lisSkinsStorage := &adapterlisskins.Storage{Client: adapterlisskins.NewLisSkinsClient(cfg), BaseURL: cfg.LisSkinsBaseURL, Logger: a.logger}
	cstmStorage := &adaptercstm.Storage{
		Client:         adaptercstm.NewClient(cfg),
		BaseURL:        cfg.CSTMBaseURL,
		RequestTimeout: cfg.CSTMRequestTimeout,
		Logger:         a.logger,
		CacheTTL:       cfg.CacheTTL,
	}
	searchNewSkinsUC := skins.SearchNewSkins{
		Storage: steamStorage,
	}
	saveSkinStorage := &adapterdbskins.Storage{
		Conn:             a.backend.Factory.DBConnection(),
		SteamStorage:     steamStorage,
		LisSkinsStorage:  lisSkinsStorage,
		CSTMStorage:      cstmStorage,
		CSTMBaseURL:      cfg.CSTMBaseURL,
		BatchUpdateDelay: cfg.BulkPriceUpdateDelay,
		Logger:           a.logger,
	}
	saveSkinUC := skins.SaveSkin{SkinSaver: saveSkinStorage}
	getSavedSkinsUC := skins.GetSavedSkins{SavedSkinsReader: saveSkinStorage}
	updateSavedSkinPriceUC := skins.UpdateSavedSkinPrice{Updater: saveSkinStorage}
	updateAllSavedSkinsPricesUC := skins.UpdateAllSavedSkinsPrices{Updater: saveSkinStorage}
	deleteSavedSkinUC := skins.DeleteSavedSkin{Deleter: saveSkinStorage}
	keyBytes, keyErr := sharedcrypto.LoadOrCreateTokenKey("SkinPrice")
	if keyErr != nil && cfg.TokenEncryptionKey != "" {
		a.logger.Warn("fallback to TOKEN_ENCRYPTION_KEY due to key file issue", logx.ErrAttrs(keyErr)...)
		keyBytes = nil
	}
	var tokenCipher *sharedcrypto.TokenCipher
	if len(keyBytes) == 32 {
		tokenCipher, err = sharedcrypto.NewTokenCipher(keyBytes)
	} else {
		tokenCipher, err = sharedcrypto.NewTokenCipherFromBase64(cfg.TokenEncryptionKey)
	}
	if err != nil {
		panic(err)
	}
	sourceStateStorage := &adapterdbsourcestate.Storage{Conn: a.backend.Factory.DBConnection()}
	appSettingsStorage := &adapterdbappsettings.Storage{Conn: a.backend.Factory.DBConnection()}
	getLisSkinsTokenUC := skins.GetLisSkinsToken{Storage: sourceStateStorage, Cipher: tokenCipher}
	lisSkinsStorage.TokenProvider = getLisSkinsTokenUC
	saveLisSkinsTokenUC := skins.SaveLisSkinsToken{Storage: sourceStateStorage, Cipher: tokenCipher}
	hasLisSkinsTokenUC := skins.HasLisSkinsToken{Storage: sourceStateStorage}
	clearLisSkinsTokenUC := skins.ClearLisSkinsToken{Storage: sourceStateStorage}
	getAppSettingsUC := appsettings.GetAppSettings{Storage: appSettingsStorage}
	saveAppSettingsUC := appsettings.SaveAppSettings{Storage: appSettingsStorage}
	a.skinsEndpoints = presenterskins.NewEndpoints(searchNewSkinsUC, saveSkinUC, getSavedSkinsUC, updateSavedSkinPriceUC, updateAllSavedSkinsPricesUC, deleteSavedSkinUC, saveLisSkinsTokenUC, hasLisSkinsTokenUC, clearLisSkinsTokenUC)
	a.settingsEndpoints = presentersettings.NewEndpoints(getAppSettingsUC, saveAppSettingsUC)
	a.logger.Info("routes registered")
}

func (a *App) SetLisSkinsToken(payload presenterskins.SetLisSkinsTokenRequest) error {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("set lisskins token requested")
	err := a.skinsEndpoints.SetLisSkinsToken(payload)
	if err != nil {
		logRouteError(logger, "set lisskins token failed", err)
		return errx.FromError(err, "invalid token")
	}
	logger.Info("set lisskins token completed")
	return nil
}

func (a *App) HasLisSkinsToken() (presenterskins.LisSkinsTokenStatusResponse, error) {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("check lisskins token status requested")
	response, err := a.skinsEndpoints.GetLisSkinsTokenStatus()
	if err != nil {
		logRouteError(logger, "check lisskins token status failed", err)
		return presenterskins.LisSkinsTokenStatusResponse{}, errx.FromError(err, "failed to check token status")
	}
	logger.Info("check lisskins token status completed", slog.Bool("has_token", response.HasToken))
	return response, nil
}

func (a *App) ClearLisSkinsToken() error {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("clear lisskins token requested")
	err := a.skinsEndpoints.ClearLisSkinsToken()
	if err != nil {
		logRouteError(logger, "clear lisskins token failed", err)
		return errx.FromError(err, "failed to clear token")
	}
	logger.Info("clear lisskins token completed")
	return nil
}

func (a *App) UpdateSavedSkinPrice(payload presenterskins.UpdateSavedSkinPriceRequest) (presenterskins.UpdateSavedSkinPriceResponse, error) {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("update saved skin price requested", slog.String("market_hash_name", payload.MarketHashName), slog.String("currency", payload.Currency))
	response, err := a.skinsEndpoints.UpdateSavedSkinPrice(payload)
	if err != nil {
		logRouteError(logger, "update saved skin price failed", err)
		return presenterskins.UpdateSavedSkinPriceResponse{}, errx.FromError(err, "failed to update skin price")
	}
	logger.Info("update saved skin price completed", slog.String("market_hash_name", payload.MarketHashName), slog.String("currency", payload.Currency))
	return response, nil
}

func (a *App) UpdateAllSavedSkinsPrices(payload presenterskins.UpdateAllSavedSkinsPricesRequest) (presenterskins.UpdateAllSavedSkinsPricesResponse, error) {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("bulk saved skins price update requested", slog.String("currency", payload.Currency))
	response, err := a.skinsEndpoints.UpdateAllSavedSkinsPrices(payload)
	if err != nil {
		logRouteError(logger, "bulk saved skins price update failed", err)
		return presenterskins.UpdateAllSavedSkinsPricesResponse{}, errx.FromError(err, "failed to update all skin prices")
	}
	logger.Info("bulk saved skins price update completed",
		slog.String("currency", payload.Currency),
		slog.Int("updated_count", response.UpdatedCount),
		slog.Int("failed_count", response.FailedCount),
	)
	return response, nil
}

func (a *App) SearchNewSkins(filter presenterskins.SearchNewSkinsFilter) (presenterskins.NewSkinsResponse, error) {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("search new skins requested",
		slog.Int("limit", filter.Limit),
		slog.Int("offset", filter.Offset),
		slog.String("sort_column", filter.SortColumn),
		slog.String("sort_dir", filter.SortDir),
		slog.Bool("has_query", filter.MarketHashName != nil && *filter.MarketHashName != ""),
	)
	response, err := a.skinsEndpoints.SearchNewSkins(filter)
	if err != nil {
		logRouteError(logger, "search new skins failed", err)
		return presenterskins.NewSkinsResponse{}, errx.FromError(err, "failed to search skins")
	}
	logger.Info("search new skins completed",
		slog.Int("items", len(response.Items)),
		slog.Int("total_count", response.TotalCount),
	)
	return response, nil
}

func (a *App) SaveSkin(payload presenterskins.SaveSkinRequest) (presenterskins.SaveSkinResponse, error) {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("save skin requested", slog.String("market_hash_name", payload.MarketHashName))
	response, err := a.skinsEndpoints.SaveSkin(payload)
	if err != nil {
		logRouteError(logger, "save skin failed", err)
		return presenterskins.SaveSkinResponse{}, errx.FromError(err, "failed to save skin")
	}
	logger.Info("save skin completed", slog.String("market_hash_name", payload.MarketHashName), slog.Bool("created", response.Created))
	return response, nil
}

func (a *App) GetSavedSkins(filter presenterskins.GetSavedSkinsFilter) (presenterskins.SavedSkinsResponse, error) {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("get saved skins requested", slog.Int("limit", filter.Limit), slog.Int("offset", filter.Offset))
	response, err := a.skinsEndpoints.GetSavedSkins(filter)
	if err != nil {
		logRouteError(logger, "get saved skins failed", err)
		return presenterskins.SavedSkinsResponse{}, errx.FromError(err, "failed to load saved skins")
	}
	logger.Info("get saved skins completed", slog.Int("items", len(response.Items)), slog.Int("total_count", response.TotalCount))
	return response, nil
}

func (a *App) DeleteSavedSkin(payload presenterskins.DeleteSavedSkinRequest) error {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("delete saved skin requested", slog.String("market_hash_name", payload.MarketHashName))
	err := a.skinsEndpoints.DeleteSavedSkin(payload)
	if err != nil {
		logRouteError(logger, "delete saved skin failed", err)
		return errx.FromError(err, "failed to delete saved skin")
	}
	logger.Info("delete saved skin completed", slog.String("market_hash_name", payload.MarketHashName))
	return nil
}

func (a *App) GetAppSettings() (presentersettings.AppSettingsResponse, error) {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("get app settings requested")
	response, err := a.settingsEndpoints.GetAppSettings()
	if err != nil {
		logRouteError(logger, "get app settings failed", err)
		return presentersettings.AppSettingsResponse{}, errx.FromError(err, "failed to load app settings")
	}
	logger.Info("get app settings completed",
		slog.String("currency", response.Currency),
		slog.Int("auto_refresh_interval_seconds", response.AutoRefreshIntervalSeconds),
	)
	return response, nil
}

func (a *App) SaveAppSettings(payload presentersettings.SaveAppSettingsRequest) error {
	logger := logx.WithComponent(a.logger, "route")
	logger.Info("save app settings requested",
		slog.String("currency", payload.Currency),
		slog.Int("auto_refresh_interval_seconds", payload.AutoRefreshIntervalSeconds),
	)
	err := a.settingsEndpoints.SaveAppSettings(payload)
	if err != nil {
		logRouteError(logger, "save app settings failed", err)
		return errx.FromError(err, "failed to save app settings")
	}
	logger.Info("save app settings completed",
		slog.String("currency", payload.Currency),
		slog.Int("auto_refresh_interval_seconds", payload.AutoRefreshIntervalSeconds),
	)
	return nil
}

type ClientLogEvent struct {
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Component string         `json:"component"`
	Context   map[string]any `json:"context"`
}

func (a *App) LogClientEvent(payload ClientLogEvent) {
	logger := logx.WithComponent(a.logger, "frontend")
	attrs := make([]any, 0, len(payload.Context)+1)
	if payload.Component != "" {
		attrs = append(attrs, slog.String("ui_component", payload.Component))
	}
	for key, value := range payload.Context {
		attrs = append(attrs, slog.Any(key, value))
	}

	switch payload.Level {
	case "debug":
		logger.Debug(payload.Message, attrs...)
	case "warn", "warning":
		logger.Warn(payload.Message, attrs...)
	case "error":
		logger.Error(payload.Message, attrs...)
	default:
		logger.Info(payload.Message, attrs...)
	}
}

func logRouteError(logger *slog.Logger, message string, err error) {
	code := errx.CodeOf(err)
	attrs := logx.ErrAttrs(err)
	if code == errx.CodeNotFound || code == errx.CodeConflict || code == errx.CodeAlreadyExists {
		logger.Warn(message, attrs...)
		return
	}
	logger.Error(message, attrs...)
}
