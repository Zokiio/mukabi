// Package config provides configuration loading functionality
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/zokiio/mukabi/service/bot"
)

// Load reads and parses the configuration file at the given path
func Load(path string) (*bot.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", path, err)
	}
	defer file.Close()

	var config *bot.Config
	if _, err := toml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file %s: %w", path, err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// validateConfig performs basic validation of the configuration
func validateConfig(cfg *bot.Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	if cfg.Bot.Token == "" {
		return fmt.Errorf("bot token is required")
	}

	if cfg.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}

	if cfg.Database.Database == "" {
		return fmt.Errorf("database name/path is required")
	}

	return nil
}
