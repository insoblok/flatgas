package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"log"
)

type DevAccount struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

// BuildGenesis generates a valid Clique genesis.json from DevAccounts.
func BuildGenesis(accounts []DevAccount, validator common.Address) []byte {
	const (
		chainID   = 12345
		blockTime = 5
		epoch     = 30000
	)

	alloc := make(map[string]map[string]string)
	for _, acc := range accounts {
		alloc[acc.Address.Hex()] = map[string]string{
			"balance": "0xfffffffffffffffffffff", // Large initial balance
		}
	}

	// extraData = 32 bytes vanity + validator address (20 bytes) + 65 bytes padding = 32 + 20 + 65 = 117 bytes
	extraData := make([]byte, 32+20+65)
	copy(extraData[32:], validator.Bytes())

	genesis := map[string]interface{}{
		"config": map[string]interface{}{
			"chainId":                 chainID,
			"homesteadBlock":          0,
			"eip150Block":             0,
			"eip155Block":             0,
			"eip158Block":             0,
			"byzantiumBlock":          0,
			"constantinopleBlock":     0,
			"petersburgBlock":         0,
			"istanbulBlock":           0,
			"berlinBlock":             0,
			"londonBlock":             0,
			"terminalTotalDifficulty": "0x0",
			"clique": map[string]interface{}{
				"period": blockTime,
				"epoch":  epoch,
			},
		},
		"nonce":      "0x0",
		"timestamp":  "0x0",
		"extraData":  "0x" + hex.EncodeToString(extraData),
		"gasLimit":   "0x47b760", // ~4700000
		"difficulty": "0x1",
		"mixHash":    "0x0000000000000000000000000000000000000000000000000000000000000000",
		"coinbase":   "0x0000000000000000000000000000000000000000",
		"alloc":      alloc,
		"number":     "0x0",
		"gasUsed":    "0x0",
		"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	}

	data, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal genesis: %v", err)
	}
	return data
}
