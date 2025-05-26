package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultRPC string            `json:"defaultRpc"`
	RPCs       map[string]string `json:"rpcs"`
}

func ConfigStore(base string) (DataStoreConfig, error) {
	// Step 1: Validate that the base path exists, is readable, and writeable
	info, err := os.Stat(base)
	if err != nil {
		return DataStoreConfig{}, fmt.Errorf("base directory does not exist or cannot be accessed: %w", err)
	}
	if !info.IsDir() {
		return DataStoreConfig{}, fmt.Errorf("base path is not a directory")
	}

	// Check if write permissions are available
	testFilePath := filepath.Join(base, ".test")
	file, err := os.Create(testFilePath)
	if err != nil {
		return DataStoreConfig{}, fmt.Errorf("base directory is not writeable: %w", err)
	}
	file.Close()
	os.Remove(testFilePath) // Clean up the temporary test file

	// Step 2: Initialize configuration values
	walletName := "wallet"
	storeBaseName := "config"
	configDir := filepath.Join(base, walletName, storeBaseName)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return DataStoreConfig{}, fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	storeFileName := "config.json"
	storeFilePath := filepath.Join(configDir, storeFileName)
	archiveDir := filepath.Join(configDir, "archive")
	if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
		if err := os.MkdirAll(archiveDir, 0700); err != nil {
			return DataStoreConfig{}, fmt.Errorf("failed to create archive directory: %w", err)
		}
	}

	// Step 3: Calculate StoreCurrentVersion based on the number of files in the archive directory
	storeCurrentVersion, err := calculateVersion(archiveDir)
	if err != nil {
		return DataStoreConfig{}, err
	}

	// Step 4: Return the populated DataStoreConfig struct
	return DataStoreConfig{
		BaseDir:             base,
		StoreBaseName:       storeBaseName,
		StoreBaseDir:        configDir,
		StoreFileName:       storeFileName,
		StoreFilePath:       storeFilePath,
		ArchiveDir:          archiveDir,
		StoreCurrentVersion: storeCurrentVersion,
	}, nil
}

// Helper function to count the number of files in the archive directory
func calculateVersion(archiveDir string) (int, error) {
	// Ensure archive directory exists or create it
	if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
		if err := os.MkdirAll(archiveDir, os.ModePerm); err != nil {
			return 0, fmt.Errorf("failed to create archive directory: %w", err)
		}
	}

	// Count files in archiveDir
	files, err := ioutil.ReadDir(archiveDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read archive directory: %w", err)
	}

	// Return count of files + 1 as the current store version
	return len(files) + 1, nil
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
