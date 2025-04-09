// Package db provides database access and operations for the bot
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

const (
	defaultTimeout = 10 * time.Second
)

// Config holds database configuration parameters
type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"ssl_mode"`
	Driver   string `toml:"driver"`
}

// String returns a string representation of the config, masking sensitive data
func (c Config) String() string {
	return fmt.Sprintf("\n   Host: %s\n   Port: %d\n   Username: %s\n   Password: %s\n   Database: %s\n   SSLMode: %s",
		c.Host,
		c.Port,
		c.Username,
		strings.Repeat("*", len(c.Password)),
		c.Database,
		c.SSLMode,
	)
}

// PostgresDataSourceName returns the PostgreSQL connection string
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

// New creates a new database connection based on the provided configuration
func New(driver string, cfg Config, schema string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	switch strings.ToLower(driver) {
	case "postgres":
		return newPostgres(ctx, cfg, schema)
	case "sqlite":
		return newSQLite(ctx, cfg, schema)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}

// Database represents a database connection with query capabilities
type Database struct {
	db *sqlx.DB
}

func newPostgres(ctx context.Context, cfg Config, schema string) (*Database, error) {
	pgCfg, err := pgx.ParseConfig(cfg.PostgresDataSourceName())
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	db, err := sqlx.Open("pgx", stdlib.RegisterConnConfig(pgCfg))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if schema != "" {
		if _, err = db.ExecContext(ctx, schema); err != nil {
			return nil, fmt.Errorf("failed to execute schema: %w", err)
		}
	}

	return &Database{db: db}, nil
}

func newSQLite(ctx context.Context, cfg Config, schema string) (*Database, error) {
	if cfg.Database == "" {
		return nil, fmt.Errorf("SQLite database path is required")
	}

	db, err := sqlx.Open("sqlite3", cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	if schema != "" {
		if _, err = db.ExecContext(ctx, schema); err != nil {
			return nil, fmt.Errorf("failed to execute schema: %w", err)
		}
	}

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
