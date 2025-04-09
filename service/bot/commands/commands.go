// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"

	"github.com/zokiio/mukabi/service/bot"
)

// Commander handles Discord slash command interactions
type Commander struct {
	*bot.Bot
}

// Commands returns all registered slash commands for the bot
func Commands() []discord.ApplicationCommandCreate {
	cmds := make([]discord.ApplicationCommandCreate, len(registry))
	for i, cmd := range registry {
		cmds[i] = cmd.Definition()
	}
	return cmds
}

// New creates a new command router with all registered commands and middlewares
func New(b *bot.Bot) handler.Router {
	cmds := &Commander{b}
	router := handler.New()
	router.Use(middleware.Go)

	// Register all commands from the registry
	for _, cmd := range registry {
		def := cmd.Definition().(discord.SlashCommandCreate)
		if handler := cmd.Handler(cmds); handler != nil {
			router.Command("/"+def.Name, handler)
		}
		if autoHandler := cmd.AutocompleteHandler(cmds); autoHandler != nil {
			router.Autocomplete("/"+def.Name, autoHandler)
		}
	}

	return router
}
