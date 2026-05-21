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
	"context"
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
	sourceStateStorage := &adapterdbsourcestate.Storage{Conn: f.dbConnection}
	priceCollector := skins.DefaultSavedSkinPriceCollector{
		Sources: []skins.PriceSource{
			skins.ReaderPriceSource{SourceID: "steam", SourceLabel: "Steam", Reader: steamStorage},
			skins.ReaderPriceSource{SourceID: "lisskins", SourceLabel: "LisSkins", Reader: lisSkinsStorage, Pages: lisSkinsStorage, CurrencyOverride: skins.LisSkinsCurrency},
			skins.ReaderPriceSource{SourceID: "cstm", SourceLabel: "CS TM", Reader: cstmStorage, Pages: cstmStorage},
		},
		SourceState: sourceStateStorage,
	}
	updateSavedSkinPriceUC := skins.UpdateSavedSkinPrice{
		Repository: savedSkinsRepository,
		Collector:  priceCollector,
	}
	refreshQueue := skins.NewRefreshQueue(updateSavedSkinPriceUC)
	refreshCtx, refreshCancel := context.WithCancel(context.Background())
	refreshQueue.Run(refreshCtx)
	f.refreshQueue = refreshQueue
	f.refreshCancel = refreshCancel
	queuedUpdateSavedSkinPriceUC := skins.QueuedSavedSkinPriceUpdater{
		Queue: refreshQueue,
		Kind:  skins.RefreshTaskManual,
	}
	updateAllSavedSkinsPricesUC := skins.UpdateAllSavedSkinsPrices{
		Repository:       savedSkinsRepository,
		UpdateOne:        queuedUpdateSavedSkinPriceUC,
		BatchUpdateDelay: cfg.BulkPriceUpdateDelay,
	}
	deleteSavedSkinUC := skins.DeleteSavedSkin{Repository: savedSkinsRepository}

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

	getSourceStatesUC := skins.ListSourceStates{Storage: sourceStateStorage}
	getDiagnosticsUC := skins.GetDiagnostics{
		SourceStates: sourceStateStorage,
		Version:      "dev",
		DatabasePath: f.dbConnection.DatabasePath(),
		LogPath:      cfg.LogFilePath,
	}

	f.skinsEndpoints = presenterskins.NewEndpoints(presenterskins.EndpointDeps{
		SearchNewSkins:            searchNewSkinsUC,
		SaveSkin:                  saveSkinUC,
		GetSavedSkins:             getSavedSkinsUC,
		UpdateSavedSkinPrice:      queuedUpdateSavedSkinPriceUC,
		UpdateAllSavedSkinsPrices: updateAllSavedSkinsPricesUC,
		DeleteSavedSkin:           deleteSavedSkinUC,
		SaveLisSkinsToken:         saveLisSkinsTokenUC,
		HasLisSkinsToken:          hasLisSkinsTokenUC,
		ClearLisSkinsToken:        clearLisSkinsTokenUC,
		GetPriceSourceStates:      getSourceStatesUC,
		GetDiagnostics:            getDiagnosticsUC,
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
