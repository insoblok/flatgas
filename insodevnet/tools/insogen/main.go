package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("🔧 Flatgas Genesis Generator")

	// Example: generate 3 dev accounts
	accounts, err := GenerateAccounts(3, "flatgas")
	if err != nil {
		log.Fatalf("Failed to generate accounts: %v", err)
	}

	// Print addresses
	for i, acc := range accounts {
		fmt.Printf("Account %d: %s\n", i+1, acc.Address.Hex())
	}

	// Use first account as validator
	validator := accounts[0].Address

	// Generate genesis file
	genesis := BuildGenesis(accounts, validator)

	// Write to file
	file := "insodevnet/genesis/genesis.json"
	err = os.WriteFile(file, genesis, 0644)
	if err != nil {
		log.Fatalf("Failed to write genesis file: %v", err)
	}

	fmt.Printf("✅ Genesis written to %s\n", file)
}
