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
					Name:         "character",
					Description:  "Name of the character",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	},
}

func (c *commands) OnRegisterCharacter(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	region := data.String("region")
	realm := data.String("realm")
	character := data.String("character")

	characterData, err := c.Bot.External.RaiderIO().FetchCharacterProfile(region, realm, character, raiderio.WithFields(
		raiderio.FieldMythicPlusScoresBySeason,
	))
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
	character := data.String("character")

	// Fetch character stats from the database
	characterData, err := c.Database.WoWGetCharacter(e.GuildID().String(), e.User().ID.String(), character)
	if err != nil {
		slog.Error("Error fetching character stats:", slog.String("error", err.Error()))
		return e.CreateMessage(res.CreateError("❌ Failed to fetch character stats. Please try again later."))
	}

	// Fetch character profile from RaiderIO
	profile, err := c.Bot.External.RaiderIO().FetchCharacterProfile(characterData.Region, characterData.Realm, character, raiderio.WithFields(
		raiderio.FieldMythicPlusScoresBySeason,
	))
	if err != nil {
		slog.Warn("Character not found",
			slog.String("region", characterData.Region),
			slog.String("realm", characterData.Realm),
			slog.String("character", character),
			tint.Err(err),
		)
		return e.CreateMessage(res.CreateError("❌ Character not found. Please double-check the spelling and try again."))
	}

	msg := &discord.MessageCreate{
		Embeds: []discord.Embed{CreateCharacterEmbed(profile)},
	}

	return e.CreateMessage(*msg)
}

func (c *commands) OnWowAutocomplete(e *handler.AutocompleteEvent) error {
	slog.Debug("Autocomplete for /wow command")
	slog.Debug(fmt.Sprintf("Focused option: %s", e.Data.Focused().Name))

	// Handle the specific subcommands and their autocomplete options
	switch *e.Data.SubCommandName {
	case "reg-character":
		if e.Data.Focused().Name == "realm" {
			return c.OnRealmAutocomplete(e)
		}
	case "char-stats":
		if e.Data.Focused().Name == "character" {
			return c.OnRegisterdCharacterAutocomplete(e)
		}
	}

	return nil
}

func (c *commands) OnRegisterdCharacterAutocomplete(e *handler.AutocompleteEvent) error {
	region := strings.ToLower(e.Data.String("region"))
	query := e.Data.String("character")
	slog.Debug("Chosen region: '%s', Autocomplete query for character: '%s'\n", region, query)

	characters, err := c.Database.WoWGetCharacters(e.GuildID().String(), e.User().ID.String())
	if err != nil {
		slog.Error("Error fetching registered characters:", slog.String("error", err.Error()))
		return nil
	}

	choices := []discord.AutocompleteChoice{}
	charactersAdded := 0

	// Return the filtered character choices based on the region and query
	for _, character := range characters {
		if strings.Contains(strings.ToLower(character.CharacterName), strings.ToLower(query)) {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  character.CharacterName,
				Value: character.CharacterName,
			})

			charactersAdded++
			if charactersAdded >= 25 {
				break
			}
		}
	}

	return e.AutocompleteResult(choices)
}

func (c *commands) OnRealmAutocomplete(e *handler.AutocompleteEvent) error {
	region := strings.ToLower(e.Data.String("region"))
	query := e.Data.String("realm")
	slog.Debug("Chosen region: '%s', Autocomplete query for realm: '%s'\n", region, query)

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

	slog.Debug(fmt.Sprintf("Filtered realms for region '%s': %v", region, choices))
	return e.AutocompleteResult(choices)
}

func CreateCharacterEmbed(character *raiderio.CharacterProfile) discord.Embed {
	embed := discord.Embed{
		Type:  "rich", // explicitly required
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
