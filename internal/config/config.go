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

const configFileName = ".gatorconfig.json"

func READ() (Config, error) {
	newStruct := Config{}

	filePath, err := getConfigFilePath()
	if err != nil {
		return newStruct, err
	}

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

func getConfigFilePath() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDirectory, configFileName), nil
}

func write(c Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return err
	}

	return nil
}