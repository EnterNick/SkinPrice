package factory

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"SkinPrice/skinprice/internal/shared/logx"
	"errors"
	"fmt"
	"log/slog"
)

type Factory struct {
	dbConnection *database.Connection
	logger       *slog.Logger
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
	logger.Info("factory initialized")
	return &Factory{
		dbConnection: connection,
		logger:       logger,
	}, nil
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

func (f *Factory) GetCurrentPrice(skinName string) (float64, error) {
	return 123.123, nil
}

func (f *Factory) DBConnection() *database.Connection {
	return f.dbConnection
}
