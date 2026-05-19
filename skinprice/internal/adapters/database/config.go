package database

import (
	"SkinPrice/skinprice/internal/shared/utils"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect"
)

type Config struct {
	Driver          string
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	Debug           bool
}

func LoadConfig() *Config {
	port, _ := strconv.Atoi(os.Getenv("APP_DB_PORT"))
	maxOpenConns, _ := strconv.Atoi(os.Getenv("APP_DB_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("APP_DB_MAX_IDLE_CONNS"))
	connMaxLifeTimeSeconds, _ := strconv.Atoi(os.Getenv("APP_DB_CONN_MAX_LIFE_TIME"))
	dbName := os.Getenv("APP_DB_NAME")
	if dbName == "" {
		dbName = os.Getenv("APP_DB_PATH")
	}
	driver := utils.GetStrWDefault("APP_DB_DRIVER", "sqlite3")
	if dbName == "" && (driver == "sqlite3" || driver == "sqlite") {
		dbName = defaultSQLiteDBPath()
	}
	return &Config{
		Host:            os.Getenv("APP_DB_HOST"),
		Port:            port,
		DBName:          dbName,
		Password:        os.Getenv("APP_DB_PASSWORD"),
		User:            os.Getenv("APP_DB_USER"),
		SSLMode:         utils.GetStrWDefault("APP_DB_SSLMODE", "disable"),
		Debug:           os.Getenv("APP_DB_DEBUG") == "true",
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		Driver:          driver,
		ConnMaxLifetime: time.Duration(connMaxLifeTimeSeconds) * time.Second,
	}
}

func defaultSQLiteDBPath() string {
	switch runtime.GOOS {
	case "windows":
		localAppData := strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
		if localAppData != "" {
			return filepath.Join(localAppData, "SkinPrice", "skinprice.db")
		}
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err == nil && homeDir != "" {
			return filepath.Join(homeDir, "Library", "Application Support", "SkinPrice", "skinprice.db")
		}
	default:
		homeDir, err := os.UserHomeDir()
		if err == nil && homeDir != "" {
			return filepath.Join(homeDir, ".local", "share", "SkinPrice", "skinprice.db")
		}
	}

	return "skinprice.db"
}

func (c Config) DSN() string {
	if c.Driver == "sqlite3" || c.Driver == "sqlite" {
		dbName := c.DBName
		if dbName == "" {
			dbName = ":memory:"
		}
		separator := "?"
		if strings.Contains(dbName, "?") {
			separator = "&"
		}
		if strings.Contains(dbName, "_fk=") {
			return dbName
		}
		return dbName + separator + "_fk=1"
	}

	var scheme string
	switch c.Driver {
	case "pgx", "postgres", "postgresql", "":
		scheme = "postgres"
	default:
		scheme = c.Driver
	}

	u := &url.URL{
		Scheme: scheme,
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port)),
		Path:   c.DBName,
	}

	q := u.Query()
	if c.SSLMode != "" {
		q.Set("sslmode", c.SSLMode)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (c Config) EntDialect() string {
	switch c.Driver {
	case "sqlite3", "sqlite":
		return dialect.SQLite
	default:
		return dialect.Postgres
	}
}
