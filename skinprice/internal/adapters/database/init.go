package database

import (
	"SkinPrice/skinprice/internal/adapters/database/ent"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // need to load db driver

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

type Connection struct {
	cfg    *Config
	sqlDB  *sql.DB
	driver dialect.Driver
	client *ent.Client
}

func NewConnection(cfg *Config, sqlDB *sql.DB, driver dialect.Driver, client *ent.Client) *Connection {
	return &Connection{
		cfg:    cfg,
		sqlDB:  sqlDB,
		driver: driver,
		client: client,
	}
}

func New(cfg *Config) (*Connection, error) {
	if cfg == nil {
		cfg = LoadConfig()
	}
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	client := ent.NewClient(ent.Driver(drv))
	if cfg.Debug {
		client = client.Debug()
	}

	return NewConnection(
		cfg,
		db,
		drv,
		client,
	), nil
}

func (a *Connection) Client() *ent.Client { return a.client }

func (a *Connection) Close() error {
	if a.client != nil {
		_ = a.client.Close()
	}
	if a.driver != nil {
		_ = a.driver.Close()
	}
	if a.sqlDB != nil {
		return a.sqlDB.Close()
	}
	return nil
}

func (a *Connection) DB() *sql.DB {
	return a.sqlDB
}
