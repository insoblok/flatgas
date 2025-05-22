package internal

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultRPC string            `json:"defaultRpc"`
	RPCs       map[string]string `json:"rpcs"`
}

// ConfigPath returns the full path to the wallet config file.
func ConfigPath(base string) string {
	return filepath.Join(base, "wallet", "config.json")
}

// LoadConfig loads the config from wallet/config.json.
func LoadConfig(base string) (Config, error) {
	path := ConfigPath(base)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return Config{RPCs: make(map[string]string)}, nil // return empty config
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.RPCs == nil {
		cfg.RPCs = make(map[string]string)
	}

	return cfg, nil
}

// SaveConfig writes the config to wallet/config.json.
func SaveConfig(base string, cfg Config) error {
	path := ConfigPath(base)

	// Ensure wallet dir exists
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
