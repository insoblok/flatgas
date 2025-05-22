package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage Flatgas dev accounts",
	Long:  "Create, list, import, and export accounts stored in the local wallet",
}

func GetAccountsCommand() *cobra.Command {
	return accountsCmd
}

var (
	createAlias    string
	createPassword string
)

var createAccountCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new account",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		if createAlias == "" || createPassword == "" {
			return fmt.Errorf("--alias and --password are required")
		}

		keystoreDir := filepath.Join(base, "wallet", "keystore")
		ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

		account, err := ks.NewAccount(createPassword)
		if err != nil {
			return fmt.Errorf("failed to create new account: %w", err)
		}

		fmt.Printf("✅ New account created: %s (alias: %s)\n", account.Address.Hex(), createAlias)
		fmt.Printf("🔐 Keyfile stored at: %s\n", account.URL.Path)

		aliasesPath := filepath.Join(base, "wallet", "aliases.json")
		var aliases map[string]string

		// Load existing aliases if the file exists
		if data, err := os.ReadFile(aliasesPath); err == nil {
			_ = json.Unmarshal(data, &aliases)
		} else {
			aliases = make(map[string]string)
		}

		// Add the new alias
		if _, exists := aliases[createAlias]; exists {
			return fmt.Errorf("alias '%s' already exists", createAlias)
		}
		aliases[createAlias] = account.Address.Hex()

		// Save back to file
		aliasData, err := json.MarshalIndent(aliases, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal aliases: %w", err)
		}
		if err := os.WriteFile(aliasesPath, aliasData, 0644); err != nil {
			return fmt.Errorf("failed to write aliases file: %w", err)
		}

		fmt.Printf("📝 Alias '%s' added to aliases.json\n", createAlias)

		return nil
	},
}

var listAccountsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all known accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		aliasesPath := filepath.Join(base, "wallet", "aliases.json")
		data, err := os.ReadFile(aliasesPath)
		if err != nil {
			return fmt.Errorf("failed to read aliases file: %w", err)
		}

		var aliases map[string]string
		if err := json.Unmarshal(data, &aliases); err != nil {
			return fmt.Errorf("failed to parse aliases: %w", err)
		}

		fmt.Println("📁 Known accounts:")
		for alias, address := range aliases {
			fmt.Printf("  %s => %s\n", alias, address)
		}
		return nil
	},
}

func init() {
	createAccountCmd.Flags().StringVar(&createAlias, "alias", "", "Alias for the new account")
	createAccountCmd.Flags().StringVar(&createPassword, "password", "", "Password to encrypt the keyfile")
	createAccountCmd.Flags().String("base", ".", "Base path to flatgas root")

	listAccountsCmd.Flags().String("base", ".", "Base path to flatgas root")

	accountsCmd.AddCommand(createAccountCmd)
	accountsCmd.AddCommand(listAccountsCmd)
}
