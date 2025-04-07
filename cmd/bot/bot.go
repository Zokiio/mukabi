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

var (
	Version = "dev"
	Commit  = "local"
)

func main() {
	cfgPath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	slog.Info("Bot is starting...", slog.String("config", *cfgPath))

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		slog.Error("failed to load config", tint.Err(err))
		return
	}

	slog.Info("Config loaded", slog.String("config", cfg.String()))
	log.Setup(cfg.Log)

	b, err := bot.New(*cfg, Version, Commit)
	if err != nil {
		slog.Error("Failed to create bot", slog.Any("error", err))
		return
	}
	defer b.Close()

	b.Discord.AddEventListeners(
		commands.New(b),
		events.New(b),
	)

	if err = b.Start(commands.Commands); err != nil {
		slog.Error("Failed to start bot", slog.Any("error", err))
		return
	}

	slog.Info("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
