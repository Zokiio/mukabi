package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/zokiio/mukabi/service/bot"
)

func Load(path string) (*bot.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}

	var config *bot.Config
	_, err = toml.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}
