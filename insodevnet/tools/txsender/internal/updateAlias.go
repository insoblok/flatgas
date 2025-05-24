package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func UpdateAlias(base string, alias string, address string) error {
	aliasesPath := filepath.Join(base, "wallet", "aliases.json")
	aliases := map[string]string{}

	// Load if exists
	if _, err := os.Stat(aliasesPath); err == nil {
		data, err := os.ReadFile(aliasesPath)
		if err != nil {
			return fmt.Errorf("failed to read aliases file: %w", err)
		}
		if err := json.Unmarshal(data, &aliases); err != nil {
			return fmt.Errorf("invalid aliases file format: %w", err)
		}
	}

	aliases[alias] = address
	out, _ := json.MarshalIndent(aliases, "", "  ")
	return os.WriteFile(aliasesPath, out, 0644)
}
