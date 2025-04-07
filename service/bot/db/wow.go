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

func (d *DB) WoWHasRegisteredCharacter(serverID, discordID string) (bool, error) {
	var count int
	err := d.dbx.DB.QueryRow(`SELECT COUNT(*) FROM wow_characters WHERE server_id = $1 AND discord_id = $2`, serverID, discordID).Scan(&count)
	if err != nil {
		slog.Error("Error checking character registration", "error", err, slog.String("serverID", serverID), slog.String("discordID", discordID))
		return false, err
	}
	return count > 0, nil
}
