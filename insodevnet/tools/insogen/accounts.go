package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func GenerateAccounts(n int, passphrase string) ([]DevAccount, error) {
	if n <= 0 {
		return nil, fmt.Errorf("account count must be positive")
	}

	dir := "insodevnet/keys/devaccounts"
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create accounts directory: %v", err)
	}

	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	var devAccounts []DevAccount

	for i := 0; i < n; i++ {
		acc, err := ks.NewAccount(passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %v", err)
		}

		// Unlock to access the private key
		err = ks.Unlock(acc, passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to unlock account: %v", err)
		}

		keyJSON, err := os.ReadFile(acc.URL.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %v", err)
		}

		key, err := keystore.DecryptKey(keyJSON, passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt key: %v", err)
		}

		devAccounts = append(devAccounts, DevAccount{
			Address:    acc.Address,
			PrivateKey: key.PrivateKey,
		})

		fmt.Printf("🔐 Account %d: %s\n", i+1, acc.Address.Hex())
	}

	return devAccounts, nil
}
