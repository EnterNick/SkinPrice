package factory

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
	"SkinPrice/skinprice/internal/shared/logx"
	sharedcrypto "SkinPrice/skinprice/internal/shared/utils/crypto"
)

func (f *Factory) buildEndpoints() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	steamStorage := &adaptersteam.Storage{Client: adaptersteam.NewSteamClient(cfg), BaseURL: cfg.SteamBaseURL, Logger: f.logger}
	lisSkinsStorage := &adapterlisskins.Storage{Client: adapterlisskins.NewLisSkinsClient(cfg), BaseURL: cfg.LisSkinsBaseURL, Logger: f.logger}
	cstmStorage := &adaptercstm.Storage{
		Client:         adaptercstm.NewClient(cfg),
		BaseURL:        cfg.CSTMBaseURL,
		RequestTimeout: cfg.CSTMRequestTimeout,
		Logger:         f.logger,
		CacheTTL:       cfg.CacheTTL,
	}

	savedSkinsRepository := &adapterdbskins.Storage{Conn: f.dbConnection}
	searchNewSkinsUC := skins.SearchNewSkins{Storage: steamStorage}
	saveSkinUC := skins.SaveSkin{Repository: savedSkinsRepository, LisSkinsPages: lisSkinsStorage, CSTMPages: cstmStorage}
	getSavedSkinsUC := skins.GetSavedSkins{Repository: savedSkinsRepository}
	priceCollector := skins.DefaultSavedSkinPriceCollector{
		SteamSource:    steamStorage,
		LisSkinsSource: lisSkinsStorage,
		CSTMSource:     cstmStorage,
		LisSkinsPages:  lisSkinsStorage,
		CSTMPages:      cstmStorage,
	}
	updateSavedSkinPriceUC := skins.UpdateSavedSkinPrice{
		Repository: savedSkinsRepository,
		Collector:  priceCollector,
	}
	updateAllSavedSkinsPricesUC := skins.UpdateAllSavedSkinsPrices{
		Repository:       savedSkinsRepository,
		UpdateOne:        updateSavedSkinPriceUC,
		BatchUpdateDelay: cfg.BulkPriceUpdateDelay,
	}
	deleteSavedSkinUC := skins.DeleteSavedSkin{Repository: savedSkinsRepository}

	sourceStateStorage := &adapterdbsourcestate.Storage{Conn: f.dbConnection}
	tokenCipher, err := f.buildTokenCipher(cfg)
	if err != nil {
		return err
	}
	getLisSkinsTokenUC := skins.GetLisSkinsToken{Storage: sourceStateStorage, Cipher: tokenCipher}
	lisSkinsStorage.TokenProvider = getLisSkinsTokenUC
	saveLisSkinsTokenUC := skins.SaveLisSkinsToken{Storage: sourceStateStorage, Cipher: tokenCipher}
	hasLisSkinsTokenUC := skins.HasLisSkinsToken{Storage: sourceStateStorage}
	clearLisSkinsTokenUC := skins.ClearLisSkinsToken{Storage: sourceStateStorage}

	appSettingsStorage := &adapterdbappsettings.Storage{Conn: f.dbConnection}
	getAppSettingsUC := appsettings.GetAppSettings{Storage: appSettingsStorage}
	saveAppSettingsUC := appsettings.SaveAppSettings{Storage: appSettingsStorage}

	f.skinsEndpoints = presenterskins.NewEndpoints(presenterskins.EndpointDeps{
		SearchNewSkins:            searchNewSkinsUC,
		SaveSkin:                  saveSkinUC,
		GetSavedSkins:             getSavedSkinsUC,
		UpdateSavedSkinPrice:      updateSavedSkinPriceUC,
		UpdateAllSavedSkinsPrices: updateAllSavedSkinsPricesUC,
		DeleteSavedSkin:           deleteSavedSkinUC,
		SaveLisSkinsToken:         saveLisSkinsTokenUC,
		HasLisSkinsToken:          hasLisSkinsTokenUC,
		ClearLisSkinsToken:        clearLisSkinsTokenUC,
	})
	f.settingsEndpoints = presentersettings.NewEndpoints(presentersettings.EndpointDeps{
		GetAppSettings:  getAppSettingsUC,
		SaveAppSettings: saveAppSettingsUC,
	})
	return nil
}

func (f *Factory) buildTokenCipher(cfg config.Config) (*sharedcrypto.TokenCipher, error) {
	keyBytes, keyErr := sharedcrypto.LoadOrCreateTokenKey("SkinPrice")
	if keyErr != nil && cfg.TokenEncryptionKey != "" {
		f.logger.Warn("fallback to TOKEN_ENCRYPTION_KEY due to key file issue", logx.ErrAttrs(keyErr)...)
		keyBytes = nil
	}
	if len(keyBytes) == 32 {
		return sharedcrypto.NewTokenCipher(keyBytes)
	}
	return sharedcrypto.NewTokenCipherFromBase64(cfg.TokenEncryptionKey)
}
