// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/handler"
	"github.com/zokiio/mukabi/service/bot/res"
)

// wrapWowMiddleware wraps a command handler with WoW-specific middleware checks
func (c *Commander) wrapWowMiddleware(next func(e *handler.CommandEvent) error) func(e *handler.CommandEvent) error {
	return func(e *handler.CommandEvent) error {
		userID := e.User().ID.String()
		guildID := e.GuildID().String()

		slog.Info("Checking character registration",
			slog.String("user", userID),
			slog.String("guild", guildID),
		)

		hasCharacter, err := c.Database.WoWHasRegisteredCharacter(guildID, userID)
		if err != nil {
			return e.CreateMessage(res.CreateError("Failed to check character registration"))
		}

		slog.Debug("Character registration check result",
			slog.String("user", userID),
			slog.String("guild", guildID),
			slog.Bool("hasCharacter", hasCharacter),
		)

		if !hasCharacter {
			return e.CreateMessage(res.CreateError("No character found. Please register a character using /wow reg-character"))
		}

		return next(e)
	}
}
