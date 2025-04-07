package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/handler"
	"github.com/zokiio/mukabi/service/bot/res"
)

func (c *commands) OnWowCommand(next func(e *handler.CommandEvent) error) func(e *handler.CommandEvent) error {
	return func(e *handler.CommandEvent) error {
		user := e.User()
		guildID := e.GuildID()

		slog.Info("Checking if user has registered a character",
			slog.String("user", user.ID.String()),
			slog.String("guild", guildID.String()),
		)
		hasCharacter, err := c.Database.WoWHasRegisteredCharacter(guildID.String(), user.ID.String())
		if err != nil {
			return e.CreateMessage(res.CreateError("Error checking character registration"))
		}

		slog.Info("Has character",
			slog.String("user", user.ID.String()),
			slog.String("guild", guildID.String()),
			slog.Bool("hasCharacter", hasCharacter),
		)

		if !hasCharacter {
			return e.CreateMessage(res.CreateError("No character found. Please register a character using /wow reg-character"))
		}

		return next(e)
	}
}
