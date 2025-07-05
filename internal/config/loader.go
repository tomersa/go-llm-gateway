package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ProviderInfo struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

var Config map[string]ProviderInfo

func LoadConfig(path string) error {
	loadedConfig, err := readConfigFromFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	Config = loadedConfig
	return nil
}

func readConfigFromFile(path string) (map[string]ProviderInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config struct {
		VirtualKeys map[string]ProviderInfo `json:"virtual_keys"`
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return config.VirtualKeys, nil
}
