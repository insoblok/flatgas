package kvaccounts

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
)

func GetExportCmd() *cobra.Command {
	var alias string
	var password string
	var outPath string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export a decrypted keyfile for an alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			dbPath := internal.GetAccountsDBFilePath(base)
			db, err := bbolt.Open(dbPath, 0600, nil)
			if err != nil {
				return fmt.Errorf("failed to open DB: %w", err)
			}
			defer db.Close()

			var record internal.AliasRecord

			err = db.View(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte("aliases"))
				if bucket == nil {
					return fmt.Errorf("aliases bucket not found")
				}

				data := bucket.Get([]byte(alias))
				if data == nil {
					return fmt.Errorf("alias not found: %s", alias)
				}
				return json.Unmarshal(data, &record)
			})
			if err != nil {
				return err
			}

			keyJSON, err := json.Marshal(record.Keystore)
			if err != nil {
				return fmt.Errorf("failed to marshal keystore: %w", err)
			}

			// Test decryption
			_, err = keystore.DecryptKey(keyJSON, password)
			if err != nil {
				return fmt.Errorf("failed to decrypt key: %w", err)
			}

			if outPath == "" {
				fmt.Printf("%s\n", keyJSON)
			} else {
				absPath := filepath.Clean(outPath)
				err = os.WriteFile(absPath, keyJSON, 0600)
				if err != nil {
					return fmt.Errorf("failed to write file: %w", err)
				}
				fmt.Printf("✅ Keyfile exported to: %s\n", absPath)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&alias, "alias", "", "Alias to export")
	cmd.Flags().StringVar(&password, "password", "", "Password to decrypt key")
	cmd.Flags().StringVar(&outPath, "out", "", "Output file (optional)")
	return cmd
}
