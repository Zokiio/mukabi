// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/zokiio/mukabi/service/bot/embeds"
)

type pingCmd struct{}

func init() {
	RegisterCommand(&pingCmd{})
}

func (c *pingCmd) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "ping",
		Description: "Check bot responsiveness",
	}
}

func (c *pingCmd) Handler(cmd *Commander) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		return e.CreateMessage(embeds.Message("Pong!"))
	}
}

func (c *pingCmd) AutocompleteHandler(_ *Commander) handler.AutocompleteHandler {
	return nil
}
