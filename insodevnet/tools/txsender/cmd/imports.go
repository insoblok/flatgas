package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
)

var importFile string
var importAlias string
var importBase string

var accountsImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a keystore file and register it with an alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		keyData, err := os.ReadFile(importFile)
		if err != nil {
			return fmt.Errorf("failed to read keyfile: %w", err)
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(keyData, &raw); err != nil {
			return fmt.Errorf("invalid keystore format: %w", err)
		}

		addr, ok := raw["address"].(string)
		if !ok || len(addr) != 40 {
			return fmt.Errorf("keystore does not contain a valid address")
		}

		// Write to wallet/keystore
		keystoreDir := filepath.Join(importBase, "wallet", "keystore")
		os.MkdirAll(keystoreDir, 0700)

		timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z")
		outName := fmt.Sprintf("UTC--%s--%s", timestamp, strings.ToLower(addr))
		outPath := filepath.Join(keystoreDir, outName)

		src, err := os.Open(importFile)
		if err != nil {
			return err
		}
		defer src.Close()
		dst, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer dst.Close()
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		// Register alias
		if importAlias != "" {
			if err := internal.UpdateAlias(importBase, importAlias, "0x"+addr); err != nil {
				return fmt.Errorf("failed to update alias: %w", err)
			}
			fmt.Printf("✅ Imported as '%s' -> 0x%s\n", importAlias, addr)
		} else {
			fmt.Printf("✅ Imported keystore for 0x%s\n", addr)
		}

		fmt.Printf("🔐 Stored at: %s\n", outPath)
		return nil
	},
}

func init() {
	accountsImportCmd.Flags().StringVar(&importFile, "file", "", "Path to keystore JSON")
	accountsImportCmd.Flags().StringVar(&importAlias, "alias", "", "Alias to register")
	accountsImportCmd.Flags().StringVar(&importBase, "base", ".", "Base path to Flatgas wallet")
	accountsImportCmd.MarkFlagRequired("file")
}

func GetAccountsImportCommand() *cobra.Command {
	return accountsImportCmd
}
