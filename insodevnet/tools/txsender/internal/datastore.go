package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DataStoreConfig struct {
	BaseDir             string
	StoreBaseName       string
	StoreBaseDir        string
	StoreFileName       string
	StoreFilePath       string
	ArchiveDir          string
	StoreCurrentVersion int
}

func GetDataStoreConfig(base string, storeBaseName string, storeFileName string) (DataStoreConfig, error) {
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
	configDir := filepath.Join(base, walletName, storeBaseName)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return DataStoreConfig{}, fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	storeFilePath := filepath.Join(configDir, storeFileName)
	archiveDir := filepath.Join(configDir, "archive")
	if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
		if err := os.MkdirAll(archiveDir, 0700); err != nil {
			return DataStoreConfig{}, fmt.Errorf("failed to create archive directory: %w", err)
		}
	}

	// Step 3: Calculate StoreCurrentVersion based on the number of files in the archive directory
	storeCurrentVersion, err := CalculateVersion(archiveDir)
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
func CalculateVersion(archiveDir string) (int, error) {
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
