package db

import (
	"log/slog"
)

type WoWCharacter struct {
	CharacterName string
	Region        string
	Realm         string
}

func (d *DB) WoWRegisterCharacter(serverID string, discordID string, character WoWCharacter) error {
	_, err := d.dbx.DB.Exec(`INSERT INTO wow_characters (server_id, discord_id, character_name, region, realm) VALUES ($1, $2, $3, $4, $5)`, serverID, discordID, character.CharacterName, character.Region, character.Realm)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) WoWGetCharacters(serverID, discordID string) ([]WoWCharacter, error) {
	var characters []WoWCharacter
	rows, err := d.dbx.DB.Query(`SELECT character_name, region, realm FROM wow_characters WHERE server_id = $1 AND discord_id = $2`, serverID, discordID)
	if err != nil {
		slog.Error("Error fetching registered characters", "error", err, slog.String("serverID", serverID), slog.String("discordID", discordID))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var character WoWCharacter
		if err := rows.Scan(&character.CharacterName, &character.Region, &character.Realm); err != nil {
			slog.Error("Error scanning character row", "error", err)
			return nil, err
		}
		characters = append(characters, character)
	}
	return characters, nil
}

func (d *DB) WoWGetCharacter(serverID, discordID, characterName string) (WoWCharacter, error) {
	var character WoWCharacter
	err := d.dbx.DB.QueryRow(`SELECT character_name, region, realm FROM wow_characters WHERE server_id = $1 AND discord_id = $2 AND character_name = $3`, serverID, discordID, characterName).Scan(&character.CharacterName, &character.Region, &character.Realm)
	if err != nil {
		slog.Error("Error fetching registered character", "error", err, slog.String("serverID", serverID), slog.String("discordID", discordID), slog.String("characterName", characterName))
		return WoWCharacter{}, err
	}
	return character, nil
}

func (d *DB) WoWHasRegisteredCharacter(serverID, discordID string) (bool, error) {
	var count int
	err := d.dbx.DB.QueryRow(`SELECT COUNT(*) FROM wow_characters WHERE server_id = $1 AND discord_id = $2`, serverID, discordID).Scan(&count)
	if err != nil {
		slog.Error("Error checking character registration", "error", err, slog.String("serverID", serverID), slog.String("discordID", discordID))
		return false, err
	}
	return count > 0, nil
}
