package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	rootPath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	fullPath := filepath.Join(rootPath, configFileName)
	input, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(input, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func write(cfg Config) error {
	rootPath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(rootPath, configFileName)

	feed, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(fullPath, feed, 0644)
	return err
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	err := write(*cfg)
	return err
}
