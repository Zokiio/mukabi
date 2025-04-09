// Package db provides database operations for World of Warcraft related data
package db

import (
	"fmt"
	"log/slog"
)

// WoWCharacter represents a World of Warcraft character in the database
type WoWCharacter struct {
	CharacterName string
	Region        string
	Realm         string
}

// WoWRegisterCharacter stores a new World of Warcraft character for a Discord user
func (d *Database) WoWRegisterCharacter(serverID, discordID string, character WoWCharacter) error {
	_, err := d.db.Exec(
		`INSERT INTO wow_characters (server_id, discord_id, character_name, region, realm) 
		VALUES ($1, $2, $3, $4, $5)`,
		serverID, discordID, character.CharacterName, character.Region, character.Realm,
	)
	if err != nil {
		return fmt.Errorf("failed to register character: %w", err)
	}
	return nil
}

// WoWGetCharacters retrieves all World of Warcraft characters registered for a Discord user
func (d *Database) WoWGetCharacters(serverID, discordID string) ([]WoWCharacter, error) {
	var characters []WoWCharacter
	rows, err := d.db.Query(
		`SELECT character_name, region, realm 
		FROM wow_characters 
		WHERE server_id = $1 AND discord_id = $2`,
		serverID, discordID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch characters: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var character WoWCharacter
		if err := rows.Scan(&character.CharacterName, &character.Region, &character.Realm); err != nil {
			slog.Error("Error scanning character row", "error", err)
			return nil, fmt.Errorf("failed to scan character row: %w", err)
		}
		characters = append(characters, character)
	}
	return characters, nil
}

// WoWGetCharacter retrieves a specific World of Warcraft character for a Discord user
func (d *Database) WoWGetCharacter(serverID, discordID, characterName string) (WoWCharacter, error) {
	var character WoWCharacter
	err := d.db.QueryRow(
		`SELECT character_name, region, realm 
		FROM wow_characters 
		WHERE server_id = $1 AND discord_id = $2 AND character_name = $3`,
		serverID, discordID, characterName,
	).Scan(&character.CharacterName, &character.Region, &character.Realm)
	if err != nil {
		return WoWCharacter{}, fmt.Errorf("failed to fetch character: %w", err)
	}
	return character, nil
}

// WoWHasRegisteredCharacter checks if a Discord user has any registered World of Warcraft characters
func (d *Database) WoWHasRegisteredCharacter(serverID, discordID string) (bool, error) {
	var count int
	err := d.db.QueryRow(
		`SELECT COUNT(*) 
		FROM wow_characters 
		WHERE server_id = $1 AND discord_id = $2`,
		serverID, discordID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check character registration: %w", err)
	}
	return count > 0, nil
}
