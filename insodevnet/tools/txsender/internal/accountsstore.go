package internal

import (
	"path/filepath"
	"time"
)

type AliasRecord struct {
	Alias    string                 `json:"alias"`
	Address  string                 `json:"address"`
	Keystore map[string]interface{} `json:"keystore"`
	Metadata map[string]interface{} `json:"meta"`
	Created  time.Time              `json:"created"`
	Updated  time.Time              `json:"updated"`
}

func GetDBFilePath(base string) string {
	return filepath.Join(base, "wallet", "kvstore", "accounts.db")
}
