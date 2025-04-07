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

var WowCommand = discord.SlashCommandCreate{
	Name:        "wow",
	Description: "Commands related to World of Warcraft features.",
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
						{
							Name:  "EU",
							Value: "eu",
						},
						{
							Name:  "US",
							Value: "us",
						},
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
			Description: "Character stats",
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionString{
					Name:        "region",
					Description: "Region of the character",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "EU",
							Value: "eu",
						},
						{
							Name:  "US",
							Value: "us",
						},
					},
				},
				&discord.ApplicationCommandOptionString{
					Name:        "character",
					Description: "Name of the character",
					Required:    true,
				},
			},
		},
	},
}

func (c *commands) OnRegisterCharacter(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	region := data.String("region")
	realm := data.String("realm")
	character := data.String("character")

	characterData, err := c.Bot.External.RaiderIO().FetchCharacterProfile(region, realm, character)
	if err != nil {
		slog.Warn("Character not found",
			slog.String("region", region),
			slog.String("realm", realm),
			slog.String("character", character),
			tint.Err(err),
		)
		return e.CreateMessage(res.CreateError("❌ Character not found. Please double-check the spelling and try again."))
	}

	if err := c.Database.WoWRegisterCharacter(e.GuildID().String(), e.User().ID.String(), db.WoWCharacter{
		CharacterName: character,
		Region:        region,
		Realm:         realm,
	}); err != nil {
		slog.Error("Error registering character:", slog.String("error", err.Error()))
		return e.CreateMessage(res.CreateError("❌ Failed to register character. Please try again later."))
	}

	msg := &discord.MessageCreate{
		Embeds: []discord.Embed{CreateCharacterEmbed(characterData)},
	}

	slog.Info("Sending embed", slog.Any("msg", msg))

	return e.CreateMessage(*msg)
}

func (c *commands) OnCharacterStats(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	// Extract the options from the command data
	region := data.String("region")
	character := data.String("character")

	// TODO: Add Autocomplete for the character names
	// TODO: Collect character data from the database
	// TODO: Fetch character stats using RaiderIO

	return e.CreateMessage(res.Create(fmt.Sprintf("Stats for character %s in region %s: [example stats here]", character, region)))
}

func (c *commands) OnWowAutocomplete(e *handler.AutocompleteEvent) error {
	slog.Debug("Autocomplete for /wow command")
	slog.Debug(fmt.Sprintf("Focused option: %s", e.Data.Focused().Name))

	// Handle the specific subcommands and their autocomplete options
	switch *e.Data.SubCommandName {
	case "reg-character":
		// We want to handle the 'region' argument under these subcommands
		if e.Data.Focused().Name == "realm" {
			return c.OnRealmAutocomplete(e)
		}
	}

	return nil
}

func (c *commands) OnRealmAutocomplete(e *handler.AutocompleteEvent) error {
	region := strings.ToLower(e.Data.String("region"))
	query := e.Data.String("realm")
	slog.Debug("Chosen region: '%s', Autocomplete query for realm: '%s'\n", region, query)

	// Fetch connected realms using RaiderIO
	realms, err := c.External.RaiderIO().FetchConnectedRealms(region, query)
	if err != nil {
		slog.Error("Error fetching connected realms:", slog.String("error", err.Error()))
		return nil
	}

	choices := []discord.AutocompleteChoice{}
	realmsAdded := 0

	// Return the filtered realm choices based on the region and query
	for _, realm := range realms {
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  realm.Realm,
			Value: realm.Slug,
		})

		realmsAdded++
		if realmsAdded >= 25 {
			break
		}
	}

	// Print filtered choices for debugging
	slog.Debug(fmt.Sprintf("Filtered realms for region '%s': %v", region, choices))
	return e.AutocompleteResult(choices)
}

func CreateCharacterEmbed(character *raiderio.WoWCharacter) discord.Embed {
	embed := discord.Embed{
		Type:  "rich", // explicitly required
		Title: character.Name,
		Color: 0x00AEEF,
		Description: fmt.Sprintf("%s %s (%s)\nSpec: **%s** (%s)\nFaction: %s\nAchievement Points: %d",
			character.Race,
			character.Class,
			character.Gender,
			character.ActiveSpecName,
			character.ActiveSpecRole,
			character.Faction,
			character.AchievementPoints,
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
