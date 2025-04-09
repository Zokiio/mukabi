// Package bot provides the core functionality for the Discord bot
package bot

import (
	"fmt"
	"strings"

	"github.com/disgoorg/snowflake/v2"
	"github.com/zokiio/mukabi/internal/log"
	"github.com/zokiio/mukabi/service/bot/db"
)

// Config holds all configuration settings for the bot
type Config struct {
	Log      log.Config     `toml:"log"`
	Bot      BotConfig      `toml:"bot"`
	Database db.Config      `toml:"database"`
	External ExternalConfig `toml:"external"`
}

// String returns a string representation of the configuration, masking sensitive data
func (c Config) String() string {
	return fmt.Sprintf("\n Log: %v\n Bot: %s\n Database: %s\n",
		c.Log,
		c.Bot,
		c.Database,
	)
}

// BotConfig holds Discord-specific configuration
type BotConfig struct {
	SyncCommands bool           `toml:"sync_commands"`
	GuildIDs     []snowflake.ID `toml:"guild_ids"`
	GatewayURL   string         `toml:"gateway_url"`
	RestURL      string         `toml:"rest_url"`
	Token        string         `toml:"token"`
}

// String returns a string representation of the bot configuration, masking sensitive data
func (c BotConfig) String() string {
	return fmt.Sprintf("\n  SyncCommands: %t\n  GuildIDs: %v\n  GatewayURL: %s\n  RestURL: %s\n  Token: %s\n",
		c.SyncCommands,
		c.GuildIDs,
		c.GatewayURL,
		c.RestURL,
		strings.Repeat("*", len(c.Token)),
	)
}

// ExternalConfig holds configuration for external services
type ExternalConfig struct {
	RaiderIOKey string `toml:"raiderio_key"`
}

// DBConfig holds database-specific configuration
type DBConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"ssl_mode"`
}
