// Package bot provides the core functionality for the Discord bot
package bot

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/topi314/tint"

	"github.com/zokiio/mukabi/external"
	"github.com/zokiio/mukabi/service/bot/db"
)

//go:embed sql/schema.sql
var schema string

// Bot represents the main bot instance with all its dependencies
type Bot struct {
	Config   Config
	Version  string
	Commit   string
	Discord  bot.Client
	Database *db.Database
	External *external.Services
}

// New creates a new bot instance with the provided configuration
func New(cfg Config, version, commit string) (*Bot, error) {
	b := &Bot{
		Config:   cfg,
		Version:  version,
		Commit:   commit,
		External: external.NewServices(cfg.External.RaiderIOKey),
	}

	// Configure gateway options
	gatewayOpts := []gateway.ConfigOpt{
		gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
	}

	// Configure shard manager
	shardOpts := []sharding.ConfigOpt{
		sharding.WithGatewayConfigOpts(gatewayOpts...),
	}

	// Add custom gateway URL if provided
	if cfg.Bot.GatewayURL != "" {
		shardOpts = []sharding.ConfigOpt{
			sharding.WithGatewayConfigOpts(append(gatewayOpts,
				gateway.WithURL(cfg.Bot.GatewayURL),
				gateway.WithCompress(false),
			)...),
			sharding.WithRateLimiter(sharding.NewNoopRateLimiter()),
		}
	}

	// Configure REST client
	restOpts := []rest.ConfigOpt{}
	if cfg.Bot.RestURL != "" {
		restOpts = append(restOpts,
			rest.WithURL(cfg.Bot.RestURL),
			rest.WithRateLimiter(rest.NewNoopRateLimiter()),
		)
	}

	// Create Discord client
	client, err := disgo.New(cfg.Bot.Token,
		bot.WithShardManagerConfigOpts(shardOpts...),
		bot.WithRestClientConfigOpts(restOpts...),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds, cache.FlagVoiceStates),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create discord client: %w", err)
	}

	// Initialize database
	database, err := db.New(cfg.Database.Driver, cfg.Database, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	b.Discord = client
	b.Database = database
	return b, nil
}

// Start initializes the bot and starts listening for Discord events
func (b *Bot) Start(commands []discord.ApplicationCommandCreate) error {
	if b.Config.Bot.SyncCommands {
		slog.Info("Syncing slash commands...")

		if err := handler.SyncCommands(b.Discord, commands, b.Config.Bot.GuildIDs); err != nil {
			return fmt.Errorf("failed to sync commands: %w", err)
		}
	}

	return b.Discord.OpenShardManager(context.Background())
}

// Close gracefully shuts down the bot
func (b *Bot) Close() {
	b.Discord.Close(context.Background())
	if err := b.Database.Close(); err != nil {
		slog.Error("Error closing database connection", tint.Err(err))
	}
}
