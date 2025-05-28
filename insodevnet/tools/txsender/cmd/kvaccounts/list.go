package kvaccounts

import (
	"encoding/json"
	"fmt"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

var kvListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases stored in the kvstore",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		dbPath := internal.GetAccountsDBFilePath(base)
		db, err := bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			return fmt.Errorf("failed to open db: %w", err)
		}
		defer db.Close()

		return db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("aliases"))
			if bucket == nil {
				fmt.Println("No aliases found.")
				return nil
			}

			fmt.Println("📁 Known aliases:")
			return bucket.ForEach(func(k, v []byte) error {
				var record internal.AliasRecord
				if err := json.Unmarshal(v, &record); err != nil {
					return err
				}
				fmt.Printf("  %s => %s\n", record.Alias, record.Address)
				return nil
			})
		})
	},
}

func GetListCmd() *cobra.Command {
	return kvListCmd
}
