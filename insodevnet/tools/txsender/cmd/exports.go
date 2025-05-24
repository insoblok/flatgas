package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
)

var exportAlias string
var exportOut string
var exportBase string

var accountsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a keystore file by alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := internal.ResolveAddressOrAlias(exportBase, exportAlias)
		if err != nil {
			return fmt.Errorf("failed to resolve alias: %w", err)
		}

		keystoreDir := filepath.Join(exportBase, "wallet", "keystore")
		files, err := os.ReadDir(keystoreDir)
		if err != nil {
			return fmt.Errorf("failed to read keystore directory: %w", err)
		}

		var keyPath string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), strings.ToLower(address[2:])) {
				keyPath = filepath.Join(keystoreDir, file.Name())
				break
			}
		}
		if keyPath == "" {
			return fmt.Errorf("no keystore file found for address %s", address)
		}

		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to read keystore file: %w", err)
		}

		if exportOut != "" {
			err := os.WriteFile(exportOut, keyData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write to %s: %w", exportOut, err)
			}
			fmt.Printf("✅ Exported to %s\n", exportOut)
		} else {
			fmt.Printf("📤 Keystore JSON for %s:\n", exportAlias)
			fmt.Println(string(keyData))
		}

		return nil
	},
}

func init() {
	accountsExportCmd.Flags().StringVar(&exportAlias, "alias", "", "Alias of the account to export")
	accountsExportCmd.Flags().StringVar(&exportOut, "out", "", "Output file path (optional)")
	accountsExportCmd.Flags().StringVar(&exportBase, "base", ".", "Base path to Flatgas wallet")
	accountsExportCmd.MarkFlagRequired("alias")
}

func GetAccountsExportCommand() *cobra.Command {
	return accountsExportCmd
}
