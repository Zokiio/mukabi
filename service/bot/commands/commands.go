package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"

	"github.com/zokiio/mukabi/service/bot"
	"github.com/zokiio/mukabi/service/bot/res"
)

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	WowCommand,
}

type commands struct {
	*bot.Bot
}

func New(b *bot.Bot) handler.Router {
	cmds := &commands{b}

	router := handler.New()
	router.Use(middleware.Go)

	router.SlashCommand("/ping", cmds.OnPing)
	router.SlashCommand("/wow", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		subcommand := data.SubCommandName
		if subcommand == nil {
			return e.CreateMessage(res.CreateError("No subcommand provided"))
		}

		switch *subcommand {
		case "reg-character":
			return cmds.OnRegisterCharacter(data, e)
		case "char-stats":
			return cmds.OnWowCommand(func(e *handler.CommandEvent) error {
				return cmds.OnCharacterStats(data, e)
			})(e)
		default:
			return e.CreateMessage(res.CreateError("Unknown WoW subcommand: %s", *subcommand))
		}
	})

	router.Autocomplete("/wow", cmds.OnWowAutocomplete)

	return router
}
