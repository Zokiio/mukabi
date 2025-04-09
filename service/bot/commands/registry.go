// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Command represents a slash command with its definition and handlers
type Command interface {
	// Definition returns the slash command creation structure
	Definition() discord.ApplicationCommandCreate // Handler returns the command handler function
	Handler(c *Commander) handler.CommandHandler
	// AutocompleteHandler returns the autocomplete handler function, if any
	AutocompleteHandler(c *Commander) handler.AutocompleteHandler
}

// registry holds all registered commands
var registry []Command

// RegisterCommand adds a command to the registry
func RegisterCommand(cmd Command) {
	registry = append(registry, cmd)
}
