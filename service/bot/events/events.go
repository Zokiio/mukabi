// Package events provides Discord event handlers for the bot
package events

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"

	mubot "github.com/zokiio/mukabi/service/bot"
)

// EventHandler manages Discord event handling
type EventHandler struct {
	*mubot.Bot
}

// OnEvent implements bot.EventListener interface
func (h *EventHandler) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.GuildJoin:
		h.handleGuildJoin(e)
	case *events.GuildLeave:
		h.handleGuildLeave(e)
	case *events.Ready:
		h.handleReady(e)
	default:
		slog.Debug("Received unhandled event type")
	}
}

// handleGuildJoin processes guild join events
func (h *EventHandler) handleGuildJoin(event *events.GuildJoin) {
	slog.Info("Bot joined guild",
		slog.String("guild_name", event.Guild.Name),
		slog.String("guild_id", event.Guild.ID.String()),
	)

	if err := h.Database.RegisterServer(event.Guild.ID.String(), event.Guild.Name); err != nil {
		slog.Error("Failed to register server in database",
			slog.String("guild_id", event.Guild.ID.String()),
			slog.String("error", err.Error()),
		)
	}
}

// handleGuildLeave processes guild leave events
func (h *EventHandler) handleGuildLeave(event *events.GuildLeave) {
	slog.Info("Bot left guild",
		slog.String("guild_id", event.GuildID.String()),
	)
}

// handleReady processes the ready event when the bot connects to Discord
func (h *EventHandler) handleReady(event *events.Ready) {
	slog.Info("Bot is ready",
		slog.String("username", event.User.Username),
		slog.String("user_id", event.User.ID.String()),
	)
}

// New creates a new event handler
func New(bot *mubot.Bot) bot.EventListener {
	return &EventHandler{bot}
}
