// Package main is the entry point for the Discord bot
package main

import (
	_ "embed"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/topi314/tint"
	"github.com/zokiio/mukabi/internal/config"
	"github.com/zokiio/mukabi/internal/log"
	"github.com/zokiio/mukabi/service/bot"
	"github.com/zokiio/mukabi/service/bot/commands"
	"github.com/zokiio/mukabi/service/bot/events"
)

// Version information set during build
var (
	Version = "dev"
	Commit  = "local"
)

func main() {
	// Parse command line flags
	cfgPath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		slog.Error("Failed to load config", tint.Err(err))
		os.Exit(1)
	}

	// Setup logging
	slog.Info("Config loaded", slog.String("config", cfg.String()))
	log.Setup(cfg.Log)

	// Initialize bot
	b, err := bot.New(*cfg, Version, Commit)
	if err != nil {
		slog.Error("Failed to create bot", tint.Err(err))
		os.Exit(1)
	}
	defer b.Close()

	// Register event handlers
	b.Discord.AddEventListeners(
		commands.New(b),
		events.New(b),
	)

	// Start bot
	if err = b.Start(commands.Commands); err != nil {
		slog.Error("Failed to start bot", tint.Err(err))
		os.Exit(1)
	}

	// Wait for shutdown signal
	slog.Info("Bot is running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
