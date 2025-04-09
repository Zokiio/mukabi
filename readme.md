# Mukabi Discord Bot

A feature-rich Discord bot built in Go, providing World of Warcraft character management and integration with external services like Raider.IO.

## Features

### World of Warcraft Integration
- Character registration and management
- Character statistics via Raider.IO integration
- Realm lookup with autocomplete
- Multi-region support (EU/US)

### Discord Features
- Modern slash command support
- Rich embeds for data display
- Ephemeral messages for error handling
- Server-specific character management

### Technical Features
- Structured logging with colored output
- Multiple database support (SQLite/PostgreSQL)
- Configurable via TOML
- Graceful shutdown handling

## Requirements

- Go 1.23.6 or later
- SQLite or PostgreSQL database
- Discord bot token
- Raider.IO API key (for WoW features)

## Project Structure

```
mukabi/
├── cmd/                # Command-line applications
│   └── bot/           # Main bot executable
├── external/          # External service integrations
│   └── raiderio/     # Raider.IO API client
├── internal/          # Private application packages
│   ├── config/       # Configuration loading
│   └── log/          # Logging setup
├── service/          # Core service implementations
│   └── bot/         # Bot service implementation
└── sql/             # Database schemas
```

## Configuration

Copy `example.toml` to `config.toml` and configure:

```toml
[log]
level = 'info'
format = 'text'
add_source = true
no_color = false

[bot]
sync_commands = true
token = "your_discord_token"

[external]
raiderio_key = "your_raiderio_key"

[database]
driver = 'sqlite'
database = 'mukabi.db'
```

See `example.toml` for all available options and documentation.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/mukabi.git
   cd mukabi
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the bot:
   ```bash
   go build -o mukabi ./cmd/bot
   ```

## Usage

1. Create and configure `config.toml`
2. Run the bot:
   ```bash
   ./mukabi
   ```

### Available Commands

- `/ping` - Check bot responsiveness
- `/wow reg-character` - Register a WoW character
- `/wow char-stats` - View character statistics

## Development

### Code Style

The project follows standard Go code style guidelines:
- Use `go fmt` and `goimports` for formatting
- Follow [Effective Go](https://go.dev/doc/effective_go) conventions
- Comprehensive documentation for exported symbols
- Structured error handling with context

### Testing

Run tests with:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
