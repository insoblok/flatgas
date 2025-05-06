package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("🔧 Flatgas Genesis Generator")

	genesisTargetDir := "insodevnet/docker/single-validator/genesis"
	genesisTargetFile := genesisTargetDir + "/genesis.json"
	accountTargetDir := genesisTargetDir + "/keys/devaccounts"

	if err := os.MkdirAll(genesisTargetDir, 0700); err != nil {
		log.Fatalf("Failed to create genesis dir: %v", err)
	}
	if err := os.MkdirAll(accountTargetDir, 0700); err != nil {
		log.Fatalf("Failed to create account dir: %v", err)
	}

	accounts, err := GenerateAccounts(accountTargetDir, 3, "flatgas")
	if err != nil {
		log.Fatalf("Failed to generate accounts: %v", err)
	}

	for i, acc := range accounts {
		fmt.Printf("Account %d: %s\n", i+1, acc.Address.Hex())
	}

	validator := accounts[0].Address

	genesis := BuildGenesis(accounts, validator)

	err = os.WriteFile(genesisTargetFile, genesis, 0644)
	if err != nil {
		log.Fatalf("Failed to write genesis file: %v", err)
	}

	fmt.Printf("✅ Genesis written to %s\n", genesisTargetFile)
}
