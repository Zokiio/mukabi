// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/topi314/tint"

	"github.com/zokiio/mukabi/external/raiderio"
	"github.com/zokiio/mukabi/service/bot/db"
	"github.com/zokiio/mukabi/service/bot/res"
)

// wowCommand defines the slash command structure for World of Warcraft features
var wowCommand = discord.SlashCommandCreate{
	Name:        "wow",
	Description: "World of Warcraft features and character management",
	Options: []discord.ApplicationCommandOption{
		&discord.ApplicationCommandOptionSubCommand{
			Name:        "reg-character",
			Description: "Register a WoW character",
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionString{
					Name:        "region",
					Description: "Region of the character",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{Name: "EU", Value: "eu"},
						{Name: "US", Value: "us"},
					},
				},
				&discord.ApplicationCommandOptionString{
					Name:         "realm",
					Description:  "Realm of the character",
					Required:     true,
					Autocomplete: true,
				},
				&discord.ApplicationCommandOptionString{
					Name:        "character",
					Description: "Name of the character",
					Required:    true,
				},
			},
		},
		&discord.ApplicationCommandOptionSubCommand{
			Name:        "char-stats",
			Description: "View character statistics",
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionString{
					Name:         "character",
					Description:  "Name of the character",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	},
}

func (c *Commander) handleRegisterCharacter(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	region := data.String("region")
	realm := data.String("realm")
	character := data.String("character")

	characterData, err := c.External.RaiderIO().FetchCharacterProfile(region, realm, character, raiderio.WithFields(
		raiderio.FieldMythicPlusScoresBySeason,
	))
	if err != nil {
		slog.Warn("Character not found",
			slog.String("region", region),
			slog.String("realm", realm),
			slog.String("character", character),
			tint.Err(err),
		)
		return e.CreateMessage(res.CreateError("Character not found. Please double-check the spelling and try again."))
	}

	if err := c.Database.WoWRegisterCharacter(e.GuildID().String(), e.User().ID.String(), db.WoWCharacter{
		CharacterName: character,
		Region:        region,
		Realm:         realm,
	}); err != nil {
		slog.Error("Failed to register character", tint.Err(err))
		return e.CreateMessage(res.CreateError("Failed to register character. Please try again later."))
	}

	msg := &discord.MessageCreate{
		Embeds: []discord.Embed{createCharacterEmbed(characterData)},
	}

	return e.CreateMessage(*msg)
}

func (c *Commander) handleCharacterStats(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	character := data.String("character")

	characterData, err := c.Database.WoWGetCharacter(e.GuildID().String(), e.User().ID.String(), character)
	if err != nil {
		slog.Error("Failed to fetch character stats", tint.Err(err))
		return e.CreateMessage(res.CreateError("Failed to fetch character stats. Please try again later."))
	}

	profile, err := c.External.RaiderIO().FetchCharacterProfile(
		characterData.Region,
		characterData.Realm,
		character,
		raiderio.WithFields(raiderio.FieldMythicPlusScoresBySeason),
	)
	if err != nil {
		slog.Warn("Character not found",
			slog.String("region", characterData.Region),
			slog.String("realm", characterData.Realm),
			slog.String("character", character),
			tint.Err(err),
		)
		return e.CreateMessage(res.CreateError("Character not found. Please check if the character still exists."))
	}

	msg := &discord.MessageCreate{
		Embeds: []discord.Embed{createCharacterEmbed(profile)},
	}

	return e.CreateMessage(*msg)
}

func (c *Commander) handleWowAutocomplete(e *handler.AutocompleteEvent) error {
	slog.Debug("Processing WoW autocomplete", slog.String("focused_option", e.Data.Focused().Name))

	switch *e.Data.SubCommandName {
	case "reg-character":
		if e.Data.Focused().Name == "realm" {
			return c.handleRealmAutocomplete(e)
		}
	case "char-stats":
		if e.Data.Focused().Name == "character" {
			return c.handleCharacterAutocomplete(e)
		}
	}

	return nil
}

func (c *Commander) handleCharacterAutocomplete(e *handler.AutocompleteEvent) error {
	query := e.Data.String("character")
	characters, err := c.Database.WoWGetCharacters(e.GuildID().String(), e.User().ID.String())
	if err != nil {
		slog.Error("Failed to fetch registered characters", tint.Err(err))
		return nil
	}

	choices := make([]discord.AutocompleteChoice, 0, 25)
	for _, character := range characters {
		if strings.Contains(strings.ToLower(character.CharacterName), strings.ToLower(query)) {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  character.CharacterName,
				Value: character.CharacterName,
			})

			if len(choices) >= 25 {
				break
			}
		}
	}

	return e.AutocompleteResult(choices)
}

func (c *Commander) handleRealmAutocomplete(e *handler.AutocompleteEvent) error {
	region := strings.ToLower(e.Data.String("region"))
	query := e.Data.String("realm")

	realms, err := c.External.RaiderIO().FetchConnectedRealms(region, query)
	if err != nil {
		slog.Error("Failed to fetch connected realms", tint.Err(err))
		return nil
	}

	choices := make([]discord.AutocompleteChoice, 0, 25)
	for _, realm := range realms {
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  realm.Realm,
			Value: realm.Slug,
		})

		if len(choices) >= 25 {
			break
		}
	}

	return e.AutocompleteResult(choices)
}

// createCharacterEmbed generates a Discord embed for a character profile
func createCharacterEmbed(character *raiderio.CharacterProfile) discord.Embed {
	embed := discord.Embed{
		Type:  discord.EmbedTypeRich,
		Title: character.Name,
		Color: 0x00AEEF,
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
