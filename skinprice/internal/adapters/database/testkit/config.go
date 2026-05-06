package testkit

import (
	"SkinPrice/skinprice/internal/adapters/database"
)

func loadConfig() *database.Config {
	return &database.Config{
		Host:            "0.0.0.0",
		Port:            5432,
		DBName:          "skinprice",
		Password:        "skinprice",
		User:            "skinprice",
		SSLMode:         "disable",
		Debug:           false,
		MaxOpenConns:    0,
		MaxIdleConns:    0,
		Driver:          "pgx",
		ConnMaxLifetime: 0,
	}
}
