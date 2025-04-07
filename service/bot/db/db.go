package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"ssl_mode"`
	Driver   string `toml:"driver"`
}

func (c Config) String() string {
	return fmt.Sprintf("\n   Host: %s\n   Port: %dn   Username: %s\n   Password: %s\n   Database: %s\n   SSLMode: %s",
		c.Host,
		c.Port,
		c.Username,
		strings.Repeat("*", len(c.Password)),
		c.Database,
		c.SSLMode,
	)
}

func (c Config) PostgresDataSourceName() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
		c.SSLMode,
	)
}

func New(driver string, cfg Config, schema string) (*DB, error) {
	switch strings.ToLower(driver) {
	case "postgres":
		pgCfg, err := pgx.ParseConfig(cfg.PostgresDataSourceName())
		if err != nil {
			return nil, err
		}

		db, err := sqlx.Open("pgx", stdlib.RegisterConnConfig(pgCfg))
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}

		if _, err = db.ExecContext(ctx, schema); err != nil {
			return nil, fmt.Errorf("failed to execute schema: %w", err)
		}

		return &DB{
			dbx: db,
		}, nil
	case "sqlite":
		fmt.Println("Using SQLite for database connection.")
		if cfg.Database == "" {
			return nil, fmt.Errorf("SQLite database path is required in the config")
		}
		db, err := sqlx.Open("sqlite3", cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
		}
		if schema != "" {
			_, err = db.ExecContext(ctx, schema)
			if err != nil {
				return nil, fmt.Errorf("failed to execute SQLite schema: %w", err)
			}
		}
		return &DB{dbx: db}, nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}

type DB struct {
	dbx *sqlx.DB
}

func (d *DB) Close() error {
	return d.dbx.Close()
}
