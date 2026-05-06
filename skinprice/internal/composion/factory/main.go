package factory

import (
	"SkinPrice/skinprice/internal/adapters/database"
	"errors"
)

type Factory struct {
	dbConnection *database.Connection
}

func NewFactory() (*Factory, error) {
	connection, err := database.New(nil)
	if err != nil {
		return nil, err
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
