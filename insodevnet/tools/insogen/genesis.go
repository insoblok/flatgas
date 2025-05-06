package main

import (
	"encoding/json"
	"math/big"

	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
)

type GenesisAlloc map[common.Address]GenesisAccount

type GenesisAccount struct {
	Balance string `json:"balance"`
}

type Genesis struct {
	Config     map[string]interface{} `json:"config"`
	Nonce      string                 `json:"nonce"`
	Timestamp  string                 `json:"timestamp"`
	ExtraData  string                 `json:"extraData"`
	GasLimit   string                 `json:"gasLimit"`
	Difficulty string                 `json:"difficulty"`
	MixHash    string                 `json:"mixHash"`
	Coinbase   string                 `json:"coinbase"`
	Alloc      GenesisAlloc           `json:"alloc"`
	Number     string                 `json:"number"`
	GasUsed    string                 `json:"gasUsed"`
	ParentHash string                 `json:"parentHash"`
}

// DevAccount holds address and key info for test accounts
type DevAccount struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

func BuildGenesis(accounts []DevAccount) []byte {
	alloc := make(GenesisAlloc)

	// Allocate 100 ETH to each dev account
	initialBalance := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)) // 100 ETH

	for _, acc := range accounts {
		alloc[acc.Address] = GenesisAccount{
			Balance: initialBalance.String(),
		}
	}

	genesis := Genesis{
		Config: map[string]interface{}{
			"chainId":             12345,
			"homesteadBlock":      0,
			"eip150Block":         0,
			"eip155Block":         0,
			"eip158Block":         0,
			"byzantiumBlock":      0,
			"constantinopleBlock": 0,
			"petersburgBlock":     0,
			"istanbulBlock":       0,
			"clique": map[string]interface{}{
				"period": 5,
				"epoch":  30000,
			},
		},
		Nonce:      "0x0",
		Timestamp:  "0x0",
		ExtraData:  "0x0",      // will need to be updated with validator info if Clique is used
		GasLimit:   "0x47b760", // 4700000
		Difficulty: "0x1",
		MixHash:    "0x0000000000000000000000000000000000000000000000000000000000000000",
		Coinbase:   "0x0000000000000000000000000000000000000000",
		Alloc:      alloc,
		Number:     "0x0",
		GasUsed:    "0x0",
		ParentHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
	}

	jsonBytes, _ := json.MarshalIndent(genesis, "", "  ")
	return jsonBytes
}
