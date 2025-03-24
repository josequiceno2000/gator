package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL string `json:"db_url"`
	CurrentUsername string `json:"current_user_name,omitempty"`
}

func READ() (Config, error) {
	newStruct := Config{}
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return newStruct, err
	}

	fileName := ".gatorconfig.json"
	filePath := filepath.Join(homeDirectory, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return newStruct, err
	}

	err = json.Unmarshal(data, &newStruct)
	if err != nil {
		return newStruct, err
	}

	return newStruct, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username
	return write(*c)
}