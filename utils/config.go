package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config represents the configuration structure
type Config struct {
	DefaultEditor string `json:"default_editor"`
}

var configFile = filepath.Join(os.Getenv("HOME"), ".smeditor_config.json")

// LoadConfig loads the configuration from the config file
func LoadConfig() (Config, error) {
	var cfg Config

	// Check if the config file exists
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return cfg, nil // Return an empty config if the file does not exist
	}

	// Read the config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, err
	}

	// Unmarshal the JSON data
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(cfg Config) error {
	// Marshal the config struct to JSON
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the config file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return err
	}

	return nil
}
