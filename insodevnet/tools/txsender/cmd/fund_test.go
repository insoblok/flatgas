package cmd

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestFundDryRun(t *testing.T) {
	tempDir := t.TempDir()
	keystoreDir := filepath.Join(tempDir, "wallet", "keystore")
	require.NoError(t, os.MkdirAll(keystoreDir, 0755))

	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount("flatgas")
	require.NoError(t, err)

	aliases := map[string]string{
		"faucet": acc.Address.Hex(),
		"test1":  "0x2222222222222222222222222222222222222222",
	}
	aliasFile := filepath.Join(tempDir, "wallet", "aliases.json")
	aliasData, _ := json.MarshalIndent(aliases, "", "  ")
	require.NoError(t, os.WriteFile(aliasFile, aliasData, 0644))

	txsenderPath := filepath.Join(".", "txsender") // assumes binary is built here

	cmd := exec.Command(txsenderPath,
		"fund",
		"--rpc", "http://localhost:8545",
		"--from", "faucet",
		"--to", "test1",
		"--amount", "0.1",
		"--password", "flatgas",
		"--base", tempDir,
	)
	output, err := cmd.CombinedOutput()
	t.Log(string(output))
	require.NoError(t, err)
	require.Contains(t, string(output), "Dry run successful")
}
