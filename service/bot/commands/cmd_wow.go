// Package commands implements Discord slash command handlers for the bot
package commands

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/topi314/tint"
	"github.com/zokiio/mukabi/external/raiderio"
	"github.com/zokiio/mukabi/service/bot/db"
	"github.com/zokiio/mukabi/service/bot/embeds"
)

type wowCmd struct{}

func init() {
	RegisterCommand(&wowCmd{})
}

func (c *wowCmd) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
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
}

func (c *wowCmd) Handler(cmd *Commander) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()
		subcommand := data.SubCommandName
		if subcommand == nil {
			return e.CreateMessage(embeds.Error("No subcommand provided"))
		}

		switch *subcommand {
		case "reg-character":
			return cmd.handleRegisterCharacter(data, e)
		case "char-stats":
			return cmd.wrapWowMiddleware(func(e *handler.CommandEvent) error {
				return cmd.handleCharacterStats(data, e)
			})(e)
		default:
			return e.CreateMessage(embeds.Error("Unknown WoW subcommand: %s", *subcommand))
		}
	}
}

func (c *wowCmd) AutocompleteHandler(cmd *Commander) handler.AutocompleteHandler {
	return func(e *handler.AutocompleteEvent) error {
		slog.Debug("Processing WoW autocomplete", slog.String("focused_option", e.Data.Focused().Name))

		switch *e.Data.SubCommandName {
		case "reg-character":
			if e.Data.Focused().Name == "realm" {
				return cmd.handleRealmAutocomplete(e)
			}
		case "char-stats":
			if e.Data.Focused().Name == "character" {
				return cmd.handleCharacterAutocomplete(e)
			}
		}

		return nil
	}
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
		return e.CreateMessage(embeds.Error("Character not found. Please double-check the spelling and try again."))
	}

	if err := c.Database.WoWRegisterCharacter(e.GuildID().String(), e.User().ID.String(), db.WoWCharacter{
		CharacterName: character,
		Region:        region,
		Realm:         realm,
	}); err != nil {
		slog.Error("Failed to register character", tint.Err(err))
		return e.CreateMessage(embeds.Error("Failed to register character. Please try again later."))
	}
	return e.CreateMessage(embeds.CharacterMessage(characterData))
}

func (c *Commander) handleCharacterStats(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	character := data.String("character")

	characterData, err := c.Database.WoWGetCharacter(e.GuildID().String(), e.User().ID.String(), character)
	if err != nil {
		slog.Error("Failed to fetch character stats", tint.Err(err))
		return e.CreateMessage(embeds.Error("Failed to fetch character stats. Please try again later."))
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
		return e.CreateMessage(embeds.Error("Character not found. Please check if the character still exists."))
	}
	return e.CreateMessage(embeds.CharacterMessage(profile))
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
