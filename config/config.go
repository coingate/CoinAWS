package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultEditor string `json:"default_editor"`
}

var configFile = filepath.Join(os.Getenv("HOME"), ".smeditor_config.json")

// LoadConfig loads the configuration from the config file
func LoadConfig() (Config, error) {
	var cfg Config

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return cfg, nil // Return an empty config if the file does not exist
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

// SaveConfig saves the configuration to the config file
func SaveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, data, 0644)
}
