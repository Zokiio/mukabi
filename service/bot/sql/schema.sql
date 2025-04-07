CREATE TABLE IF NOT EXISTS servers (
    server_id TEXT PRIMARY KEY,
    server_name TEXT
);

-- Register one or more world of warcraft character for a discord user 
CREATE TABLE IF NOT EXISTS wow_characters (
    discord_id TEXT,
    server_id TEXT,
    character_name TEXT,
    region TEXT,
    realm TEXT,
    PRIMARY KEY (discord_id, server_id, character_name),
    FOREIGN KEY (server_id) REFERENCES servers(server_id)
);

