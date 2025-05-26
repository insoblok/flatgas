package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultRPC string            `json:"defaultRpc"`
	RPCs       map[string]string `json:"rpcs"`
}

func ConfigStore(base string) (DataStoreConfig, error) {
	return GetDataStoreConfig(base, "config", "config.json")
}

// LoadConfig loads the config from wallet/config.json.
func LoadConfig(base string) (Config, error) {
	config, err := ConfigStore(base)
	if err != nil {
		return Config{}, fmt.Errorf("failed to initialize config store: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(config.StoreFilePath); os.IsNotExist(err) {
		// Initialize empty config if file doesn't exist
		cfg := Config{
			DefaultRPC: "",
			RPCs:       make(map[string]string),
		}
		err := SaveConfig(base, cfg)
		if err != nil {
			return Config{}, err
		}
	}

	data, err := os.ReadFile(config.StoreFilePath)
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

func SaveConfig(base string, cfg Config) error {
	dataStoreConfig, err := ConfigStore(base)
	if err != nil {
		return fmt.Errorf("failed to initialize config store: %w", err)
	}

	archivePath := filepath.Join(dataStoreConfig.ArchiveDir, fmt.Sprintf("%s.%d", filepath.Base(dataStoreConfig.StoreFilePath), dataStoreConfig.StoreCurrentVersion))
	if err := os.Rename(dataStoreConfig.StoreFilePath, archivePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to archive existing config: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(dataStoreConfig.StoreFilePath, data, 0600)
}

func Rollback(base string) error {
	dataStoreConfig, err := ConfigStore(base)
	if err != nil {
		return fmt.Errorf("failed to initialize config store: %w", err)
	}
	if err := os.Remove(dataStoreConfig.StoreFilePath); err != nil {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	archivePath := filepath.Join(dataStoreConfig.ArchiveDir, fmt.Sprintf("%s.%d", filepath.Base(dataStoreConfig.StoreFilePath), dataStoreConfig.StoreCurrentVersion-1))
	if err := os.Rename(archivePath, dataStoreConfig.StoreFilePath); err != nil {
		return fmt.Errorf("failed to rollback config: %w", err)
	}
	return nil
}
