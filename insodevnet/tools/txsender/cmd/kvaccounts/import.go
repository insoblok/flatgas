package kvaccounts

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"io/ioutil"
	"time"
)

func GetImportCmd() *cobra.Command {
	var alias, password, keyfilePath string
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a keyfile into the kvaccounts store",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			dbPath := internal.GetAccountsDBFilePath(base)

			fmt.Printf("📁 Importing keyfile for alias '%s'...", alias)
			keyJSON, err := ioutil.ReadFile(keyfilePath)
			if err != nil {
				return fmt.Errorf("failed to read keyfile: %w", err)
			}

			key, err := keystore.DecryptKey(keyJSON, password)
			if err != nil {
				return fmt.Errorf("failed to decrypt keyfile: %w", err)
			}

			var keystoreMap map[string]interface{}
			if err := json.Unmarshal(keyJSON, &keystoreMap); err != nil {
				return fmt.Errorf("failed to decode keyfile JSON: %w", err)
			}

			record := internal.AliasRecord{
				Alias:    internal.Alias(alias),
				Address:  internal.KvAddress(key.Address.Hex()),
				Keystore: keystoreMap,
				Metadata: map[string]string{"tag": "imported"},
				Created:  time.Now(),
				Updated:  time.Now(),
			}

			db, err := bbolt.Open(dbPath, 0600, nil)
			if err != nil {
				return fmt.Errorf("failed to open db: %w", err)
			}
			defer db.Close()

			return db.Update(func(tx *bbolt.Tx) error {
				if err := internal.SaveAliasRecord(
					tx,
					record.Alias.String(),
					record,
					internal.ActionCreate,
				); err != nil {
					return fmt.Errorf("failed to store alias: %w", err)
				}

				fmt.Printf("✅ Imported alias '%s' → %s\n", alias, record.Address)
				return nil
			})

		},
	}

	cmd.Flags().StringVar(&alias, "alias", "", "Alias for the imported account")
	cmd.Flags().StringVar(&password, "password", "", "Password to decrypt the keyfile")
	cmd.Flags().StringVar(&keyfilePath, "keyfile", "", "Path to the keyfile JSON")
	cmd.MarkFlagRequired("alias")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("keyfile")

	return cmd
}
