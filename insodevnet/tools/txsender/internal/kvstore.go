package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetDBFilePathForStore(base string, storeDir string, dbFile string) string {
	dir := filepath.Join(base, "wallet", storeDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create db directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(dir, dbFile)
}

func GetAccountsDBFilePath(base string) string {
	return GetDBFilePathForStore(base, "accounts", "accounts.db")
}

func GetConfigStoreFilePath(base string) string {
	return GetDBFilePathForStore(base, "config", "config.db")
}
