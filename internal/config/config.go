package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DBURL string `json:"db_url"`
}

func Read() (Config, error) {

	dir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	fullPath := filepath.Join(dir, ".gatorconfig.json")

	readPath, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	json
}
