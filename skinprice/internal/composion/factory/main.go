package factory

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"context"
	"errors"
	"fmt"
)

type Factory struct {
	dbConnection *database.Connection
}

func NewFactory() (*Factory, error) {
	connection, err := database.New(nil)
	if err != nil {
		return nil, err
	}
	if err := connection.Client().Schema.Create(context.Background()); err != nil {
		_ = connection.Close()
		return nil, fmt.Errorf("create database schema: %w", err)
	}
	return &Factory{
		dbConnection: connection,
	}, nil
}

func (f *Factory) Close() error {
	var closeErr error

	if f.dbConnection != nil {
		if err := f.dbConnection.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
	}
	return closeErr
}

func (f *Factory) GetCurrentPrice(skinName string) (float64, error) {
	return 123.123, nil
}

func (f *Factory) DBConnection() *database.Connection {
	return f.dbConnection
}
