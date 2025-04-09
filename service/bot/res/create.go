// Package res provides Discord message response utilities
package res

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

// Color constants for message embeds
const (
	ColorPrimary = 0x5c5fea // Primary color for standard messages
	ColorDanger  = 0xd43535 // Danger color for error messages
)

// Create returns a new MessageCreate with the given content in an embed
func Create(content string) discord.MessageCreate {
	return discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Description: content,
				Color:       ColorPrimary,
			},
		},
	}
}

// Createf returns a new MessageCreate with formatted content in an embed
func Createf(format string, a ...any) discord.MessageCreate {
	return Create(fmt.Sprintf(format, a...))
}

// CreateErr returns a new error MessageCreate with the provided message and error
func CreateErr(message string, err error) discord.MessageCreate {
	return CreateError(message + ": " + err.Error())
}

// CreateError returns a new error MessageCreate with the provided message
func CreateError(message string, a ...any) discord.MessageCreate {
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
