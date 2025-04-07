# Mukabi Bot

Mukabi is a Discord bot designed to provide various features, including World of Warcraft character management and integration with external APIs like Raider.IO.

## Features

- **World of Warcraft Integration**:
  - Register WoW characters.
  - Fetch character stats using Raider.IO.
  - Autocomplete for realms and regions.

- **Discord Commands**:
  - `/ping`: Check if the bot is online.
  - `/wow reg-character`: Register a WoW character.
  - `/wow char-stats`: Fetch stats for a registered WoW character.

- **Database Support**:
  - SQLite and PostgreSQL support.
  - Schema for managing servers and WoW characters.

## Setup

### Prerequisites

- Go 1.23.6 or later
- SQLite or PostgreSQL
- A Discord bot token

### Configuration

Create a `config.toml` file in the root directory. Use the following template:

```toml
[log]
level = 'info'
format = 'text'
add_source = true
no_color = false

[bot]
dev_mode = false
sync_commands = true
guild_ids = [120041001107587072]
token = "your_discord_bot_token"

[external_api]
raiderio_key = "your_raiderio_api_key"

[database]
driver = 'sqlite'        # or 'postgres'
database = 'mukabi.db'   # SQLite database file name or PostgreSQL database name
# Uncomment the following lines for PostgreSQL
# host = 'localhost'
# port = 5432
# username = 'username'
# password = 'password'
# ssl_mode = 'disable'
```

### Running the Bot

Install dependencies:

1. Install dependencies:

```shell
go mod tidy
```

2. Run the bot:

```shell
go run ./cmd/bot
```

### Commands

`/ping`
Responds with "Pong!" to check if the bot is online.

`/wow reg-character`
Registers a World of Warcraft character. Requires the following options:

- region: The region of the character (e.g., EU, US).
- realm: The realm of the character.
- character: The name of the character.

`/wow char-stats` Fetches stats for a registered WoW character.

Requires:

- character: The name of the character.

### Logging

The bot uses a customizable logging system. Configure the log level, format, and other options in the `config.toml` file.

### External APIs

Raider.IO: Used for fetching WoW character stats and realm information.

### Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
