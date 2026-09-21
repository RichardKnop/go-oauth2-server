package config

import (
	"encoding/json"
	"os"
)

var (
	configFilePath = "/config/config.json"
)

type fileBackend struct{}

func (b *fileBackend) InitConfigBackend() {
	if os.Getenv("CONFIG_FILE_PATH") != "" {
		configFilePath = os.Getenv("CONFIG_FILE_PATH")
	}
}

func (b *fileBackend) LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (b *fileBackend) RefreshConfig(newCnf *Config) {
	*Cnf = *newCnf
}
