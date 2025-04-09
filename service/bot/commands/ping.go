// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/zokiio/mukabi/service/bot/res"
)

// pingCommand defines a simple ping-pong command for testing bot connectivity
var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Check bot responsiveness",
}

// handlePing responds to the ping command with a simple pong message
func (c *Commander) handlePing(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(res.Create("Pong!"))
}
