// Package embeds provides Discord embed creation utilities
package embeds

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

// Message creates a basic message embed
func Message(content string) discord.MessageCreate {
	return discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Description: content,
				Color:       ColorPrimary,
			},
		},
	}
}

// Messagef creates a formatted message embed
func Messagef(format string, a ...any) discord.MessageCreate {
	return Message(fmt.Sprintf(format, a...))
}

// Error creates an error message embed
func Error(message string, a ...any) discord.MessageCreate {
	return discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Description: fmt.Sprintf(message, a...),
				Color:       ColorDanger,
			},
		},
		Flags: discord.MessageFlagEphemeral,
	}
}

// ErrorWithErr creates an error message embed with error details
func ErrorWithErr(message string, err error) discord.MessageCreate {
	return Error(message + ": " + err.Error())
}
