-- Database schema for the Discord bot
-- Contains tables for server management and World of Warcraft character tracking

-- Servers table stores basic information about Discord servers the bot is in
CREATE TABLE IF NOT EXISTS servers (
    server_id TEXT PRIMARY KEY,  -- Discord server/guild ID
    server_name TEXT             -- Discord server/guild name
);

-- WoW characters table stores World of Warcraft character information for Discord users
CREATE TABLE IF NOT EXISTS wow_characters (
    discord_id TEXT,     -- Discord user ID
    server_id TEXT,      -- Discord server/guild ID
    character_name TEXT, -- WoW character name
    region TEXT,         -- WoW region (e.g., 'eu', 'us')
    realm TEXT,          -- WoW realm name
    PRIMARY KEY (discord_id, server_id, character_name),
    FOREIGN KEY (server_id) REFERENCES servers(server_id)
);

