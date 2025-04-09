// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"

	"github.com/zokiio/mukabi/service/bot"
	"github.com/zokiio/mukabi/service/bot/res"
)

// Commands contains all registered slash commands for the bot
var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	wowCommand,
}

// Commander handles Discord slash command interactions
type Commander struct {
	*bot.Bot
}

// New creates a new command router with all registered commands and middlewares
func New(b *bot.Bot) handler.Router {
	cmds := &Commander{b}

	router := handler.New()
	router.Use(middleware.Go)

	router.SlashCommand("/ping", cmds.handlePing)
	router.SlashCommand("/wow", cmds.handleWowCommand)
	router.Autocomplete("/wow", cmds.handleWowAutocomplete)

	return router
}

func (c *Commander) handleWowCommand(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	subcommand := data.SubCommandName
	if subcommand == nil {
		return e.CreateMessage(res.CreateError("No subcommand provided"))
	}

	switch *subcommand {
	case "reg-character":
		return c.handleRegisterCharacter(data, e)
	case "char-stats":
		return c.wrapWowMiddleware(func(e *handler.CommandEvent) error {
			return c.handleCharacterStats(data, e)
		})(e)
	default:
		return e.CreateMessage(res.CreateError("Unknown WoW subcommand: %s", *subcommand))
	}
}
