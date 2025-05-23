package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
)

func ResolveAddressOrAlias(base, input string) (string, error) {
	// Try alias
	aliasesPath := filepath.Join(base, "wallet", "aliases.json")
	data, err := os.ReadFile(aliasesPath)
	if err == nil {
		var aliases map[string]string
		if err := json.Unmarshal(data, &aliases); err == nil {
			if addr, ok := aliases[input]; ok {
				return addr, nil
			}
		}
	}

	// Try address
	if common.IsHexAddress(input) {
		return common.HexToAddress(input).Hex(), nil
	}

	return "", fmt.Errorf("not a valid alias or address: %s", input)
}
