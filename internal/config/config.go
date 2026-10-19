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

func (c *Config) SetUser(username string) error {

	c.CurrentUserName = username

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	path, _ := getConfigFilePath()

	err = os.WriteFile(path, data, 0600)
	if err != nil {
		return err
	}

	return write(*c)
}

func write(cfg Config) error {

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func getConfigFilePath() (string, error) {

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(dir, configFileName)

	return fullPath, nil

}

func Read() (Config, error) {

	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	readFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = json.Unmarshal(readFile, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
