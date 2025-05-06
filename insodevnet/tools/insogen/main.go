package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("🔧 Flatgas Genesis Generator")

	// Example: generate 3 dev accounts
	accounts, _ := GenerateAccounts(3, "password")

	// Print addresses
	for i, acc := range accounts {
		fmt.Printf("Account %d: %s\n", i+1, acc.Address.Hex())
	}

	// Generate genesis file
	genesis := BuildGenesis(accounts)

	// Write to file
	file := "insodevnet/genesis.json"
	err := os.WriteFile(file, genesis, 0644)
	if err != nil {
		log.Fatalf("Failed to write genesis file: %v", err)
	}

	fmt.Printf("✅ Genesis written to %s\n", file)
}
