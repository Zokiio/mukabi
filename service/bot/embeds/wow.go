// Package embeds provides Discord embed creation utilities
package embeds

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/zokiio/mukabi/external/raiderio"
)

// Character creates an embed for a WoW character profile
func Character(character *raiderio.CharacterProfile) discord.Embed {
	embed := discord.Embed{
		Type:  discord.EmbedTypeRich,
		Title: character.Name,
		Color: ColorWoW,
		Description: fmt.Sprintf(
			"**Region:** %s\n**Realm:** %s\n**Faction:** %s\n**Class:** %s\n**Mythic+ Score:** %.2f\n**Raider.IO Profile:** [Link](%s)",
			strings.ToUpper(character.Region),
			strings.ToUpper(character.Realm),
			strings.ToUpper(character.Faction),
			strings.ToUpper(character.Class),
			character.MythicPlusScoresBySeason[0].Scores.All,
			character.ProfileURL,
		),
	}

	if character.ProfileURL != "" {
		embed.URL = character.ProfileURL
	}

	if character.ThumbnailURL != "" {
		embed.Thumbnail = &discord.EmbedResource{
			URL: character.ThumbnailURL,
		}
	}

	if character.ProfileBanner != "" {
		embed.Image = &discord.EmbedResource{
			URL: character.ProfileBanner,
		}
	}

	if !character.LastCrawledAt.IsZero() {
		embed.Footer = &discord.EmbedFooter{
			Text: "Last crawled at " + character.LastCrawledAt.Format("2006-01-02 15:04:05"),
		}
	}

	return embed
}

// CharacterMessage creates a message embed for a WoW character profile
func CharacterMessage(character *raiderio.CharacterProfile) discord.MessageCreate {
	return discord.MessageCreate{
		Embeds: []discord.Embed{Character(character)},
	}
}
