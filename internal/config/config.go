package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
)

type Config struct {
	TenantID   string `json:"tenant_id"`
	ClientID   string `json:"client_id"`
	Endpoint   string `json:"endpoint"`
	Deployment string `json:"deployment"`
	APIVersion string `json:"api_version"`
}

var configDir = filepath.Join(xdg.ConfigHome, "th")
var configPath = filepath.Join(configDir, "config.json")

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found")
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

func ConfigPath() string {
	return configPath
}

func ConfigDir() string {
	return configDir
}

func IsConfigNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "config not found")
}
