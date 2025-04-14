package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	CurrentUsername string `json:"current_user_name"`
	DatabaseURL     string `json:"db_url"`
}

func Read() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	defer file.Close()
	var config Config
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}
	return &config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUsername = user
	if err := write(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) GetUser() string {
	return c.CurrentUsername
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error fetching home directory: %w", err)
	}
	path := filepath.Join(homedir, configFileName)
	return path, nil
}

func write(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling: %w", err)
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}
