package bot

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/snowflake/v2"
	"github.com/topi314/tint"
	"github.com/zokiio/mukabi/external"
	"github.com/zokiio/mukabi/service/bot/db"
)

//go:embed sql/schema.sql
var schema string

type Bot struct {
	Config   Config
	Version  string
	Commit   string
	Discord  bot.Client
	Database *db.DB
	External *external.Services
}

func New(cfg Config, version string, commit string) (*Bot, error) {
	b := &Bot{
		Config:   cfg,
		Version:  version,
		Commit:   commit,
		External: external.NewServices(cfg.External.RaiderIOKey),
	}

	gatewayConfigOpts := []gateway.ConfigOpt{
		gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
	}
	shardManagerConfigOpts := []sharding.ConfigOpt{
		sharding.WithGatewayConfigOpts(gatewayConfigOpts...),
	}
	if cfg.Bot.GatewayURL != "" {
		shardManagerConfigOpts = []sharding.ConfigOpt{
			sharding.WithGatewayConfigOpts(append(gatewayConfigOpts,
				gateway.WithURL(cfg.Bot.GatewayURL),
				gateway.WithCompress(false),
			)...),
			sharding.WithRateLimiter(sharding.NewNoopRateLimiter()),
		}
	}

	var restClientConfigOpts []rest.ConfigOpt
	if cfg.Bot.RestURL != "" {
		restClientConfigOpts = []rest.ConfigOpt{
			rest.WithURL(cfg.Bot.RestURL),
			rest.WithRateLimiter(rest.NewNoopRateLimiter()),
		}
	}

	d, err := disgo.New(cfg.Bot.Token,
		bot.WithShardManagerConfigOpts(shardManagerConfigOpts...),
		bot.WithRestClientConfigOpts(restClientConfigOpts...),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds, cache.FlagVoiceStates),
		),
	)
	if err != nil {
		slog.Error("Failed to create discord client", tint.Err(err))
		return nil, err
	}

	database, err := db.New(cfg.Database.Driver, cfg.Database, schema)
	if err != nil {
		slog.Error("Failed to create database", tint.Err(err))
		return nil, err
	}

	b.Discord = d
	b.Database = database
	return b, nil
}

func (b *Bot) Start(commands []discord.ApplicationCommandCreate) error {
	if b.Config.Bot.SyncCommands {
		slog.Info("Syncing commands...")

		// Clear existing commands
		var emptyCommands []discord.ApplicationCommandCreate
		var guildIDs []snowflake.ID
		if err := handler.SyncCommands(b.Discord, emptyCommands, guildIDs); err != nil {
			slog.Error("Failed to sync commands", tint.Err(err))
			return err
		}

		// Sync new commands
		slog.Debug("Commands to sync", "commands", commands)
		if err := handler.SyncCommands(b.Discord, commands, b.Config.Bot.GuildIDs); err != nil {
			slog.Error("Failed to sync commands", tint.Err(err))
			return err
		}
	}

	return b.Discord.OpenShardManager(context.Background())
}

func (b *Bot) Close() {
	b.Discord.Close(context.Background())
	_ = b.Database.Close()
}
