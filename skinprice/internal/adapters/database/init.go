package database

import (
	"SkinPrice/skinprice/internal/adapters/database/ent"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib" // need to load db driver
	_ "github.com/mattn/go-sqlite3"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"SkinPrice/skinprice/internal/shared/logx"
)

type Connection struct {
	cfg    *Config
	sqlDB  *sql.DB
	driver dialect.Driver
	client *ent.Client
	logger *slog.Logger
}

func NewConnection(cfg *Config, sqlDB *sql.DB, driver dialect.Driver, client *ent.Client, logger *slog.Logger) *Connection {
	return &Connection{
		cfg:    cfg,
		sqlDB:  sqlDB,
		driver: driver,
		client: client,
		logger: logger,
	}
}

func New(cfg *Config, logger ...*slog.Logger) (*Connection, error) {
	var baseLogger *slog.Logger
	if len(logger) > 0 {
		baseLogger = logger[0]
	}
	baseLogger = logx.WithComponent(baseLogger, "database")
	if cfg == nil {
		cfg = LoadConfig()
	}
	baseLogger.Info("opening database connection",
		slog.String("driver", cfg.Driver),
		slog.String("dialect", cfg.EntDialect()),
		slog.Bool("debug", cfg.Debug),
	)
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		baseLogger.Error("sql open failed", logx.ErrAttrs(err)...)
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

	drv := entsql.OpenDB(cfg.EntDialect(), db)

	client := ent.NewClient(ent.Driver(drv))
	if cfg.Debug {
		client = client.Debug()
	}

	connection := NewConnection(
		cfg,
		db,
		drv,
		client,
		baseLogger,
	)
	baseLogger.Info("database connection ready")
	return connection, nil
}

func (a *Connection) Client() *ent.Client { return a.client }

func (a *Connection) Close() error {
	logger := logx.Safe(a.logger)
	if a.client != nil {
		_ = a.client.Close()
	}
	if a.driver != nil {
		_ = a.driver.Close()
	}
	if a.sqlDB != nil {
		if err := a.sqlDB.Close(); err != nil {
			logger.Error("failed to close sql database", logx.ErrAttrs(err)...)
			return err
		}
	}
	logger.Info("database connection closed")
	return nil
}

func (a *Connection) DB() *sql.DB {
	return a.sqlDB
}

func (a *Connection) Dialect() string {
	return a.cfg.EntDialect()
}
