package events

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"

	mubot "github.com/zokiio/mukabi/service/bot"
)

type Listener struct {
	*mubot.Bot
}

// OnEvent implements bot.EventListener.
func (l *Listener) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.GuildJoin:
		l.OnGuildJoin(e)
	case *events.GuildLeave:
		l.OnGuildLeave(e)
	case *events.Ready:
		l.OnReady(e) // You'll likely want to handle the Ready event too!
	// Add cases for other event types you want to handle
	default:
		slog.Debug("Unhandled event")
	}
}

// This implements disgo's bot.EventListener
func (l *Listener) OnGuildJoin(event *events.GuildJoin) {
	slog.Info("✅ Joined guild", slog.String("guild", event.Guild.Name))
}

func (l *Listener) OnGuildLeave(event *events.GuildLeave) {
	slog.Info("❌ Left guild", slog.String("guild_id", event.GuildID.String()))
}

func (l *Listener) OnReady(event *events.Ready) {
	slog.Info("🤖 Bot is ready!")
	// Perform actions that should happen when the bot is online, like setting presence.
}

// New returns a disgo bot.EventListener
func New(bot *mubot.Bot) bot.EventListener {
	return &Listener{bot}
}
