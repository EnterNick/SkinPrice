package factory

import (
	"SkinPrice/skinprice/internal/adapters/database"
	presentersettings "SkinPrice/skinprice/internal/presenters/settings"
	presenterskins "SkinPrice/skinprice/internal/presenters/skins"
	"SkinPrice/skinprice/internal/shared/logx"
	"errors"
	"fmt"
	"log/slog"
)

type Factory struct {
	dbConnection      *database.Connection
	skinsEndpoints    *presenterskins.Endpoints
	settingsEndpoints *presentersettings.Endpoints
	logger            *slog.Logger
}

func NewFactory(logger *slog.Logger) (*Factory, error) {
	logger = logx.WithComponent(logger, "factory")
	connection, err := database.New(nil, logger)
	if err != nil {
		logger.Error("failed to initialize database connection", logx.ErrAttrs(err)...)
		return nil, err
	}
	if err := database.EnsureSchema(connection); err != nil {
		_ = connection.Close()
		logger.Error("failed to ensure database schema", logx.ErrAttrs(err)...)
		return nil, fmt.Errorf("ensure database schema: %w", err)
	}
	f := &Factory{
		dbConnection: connection,
		logger:       logger,
	}
	if err := f.buildEndpoints(); err != nil {
		_ = connection.Close()
		return nil, err
	}
	logger.Info("factory initialized")
	return f, nil
}

func (f *Factory) Close() error {
	var closeErr error

	if f.dbConnection != nil {
		if err := f.dbConnection.Close(); err != nil {
			f.logger.Error("failed to close database connection", logx.ErrAttrs(err)...)
			closeErr = errors.Join(closeErr, err)
		}
	}
	if closeErr == nil {
		f.logger.Info("factory closed")
	}
	return closeErr
}

func (f *Factory) DBConnection() *database.Connection {
	return f.dbConnection
}

func (f *Factory) SkinsEndpoints() *presenterskins.Endpoints {
	return f.skinsEndpoints
}

func (f *Factory) SettingsEndpoints() *presentersettings.Endpoints {
	return f.settingsEndpoints
}
