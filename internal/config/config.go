package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ListenAddress string `json:"listenAddress"`
	Port          int    `json:"port"`
	AuthToken     string `json:"authToken"`
	DataDir       string `json:"dataDir"`
	LogLevel      string `json:"logLevel"`
}

func LoadConfig(path string) (Config, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
