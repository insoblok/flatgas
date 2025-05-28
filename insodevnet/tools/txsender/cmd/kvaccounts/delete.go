package kvaccounts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

func GetDeleteCmd() *cobra.Command {
	var alias string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an alias from the kvstore (rollbackable)",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			if alias == "" {
				return fmt.Errorf("--alias is required")
			}

			dbPath := internal.GetAccountsDBFilePath(base)
			db, err := bbolt.Open(dbPath, 0600, nil)
			if err != nil {
				return fmt.Errorf("failed to open db: %w", err)
			}
			defer db.Close()

			var deletedRecord *internal.AliasRecord

			err = db.Update(func(tx *bbolt.Tx) error {
				aliases := tx.Bucket([]byte("aliases"))
				if aliases == nil {
					return fmt.Errorf("aliases bucket not found")
				}

				v := aliases.Get([]byte(alias))
				if v == nil {
					return fmt.Errorf("alias '%s' not found", alias)
				}

				if err := aliases.Delete([]byte(alias)); err != nil {
					return err
				}

				deletedRecord = &internal.AliasRecord{}
				if err := json.Unmarshal(v, deletedRecord); err != nil {
					return fmt.Errorf("failed to unmarshal deleted record: %w", err)
				}

				return nil
			})
			if err != nil {
				return err
			}

			entry := internal.JournalEntry{
				Action:    internal.ActionDelete,
				Alias:     internal.Alias(alias),
				Timestamp: time.Now(),
				Data:      deletedRecord,
			}

			if err := internal.WriteJournalEntry(db, entry); err != nil {
				return fmt.Errorf("failed to write journal entry: %w", err)
			}

			if err := internal.WriteAuditLogEntry(db, entry); err != nil {
				return fmt.Errorf("failed to write audit log: %w", err)
			}

			fmt.Printf("🗑️ Deleted alias '%s'\n", alias)
			return nil
		},
	}

	cmd.Flags().StringVar(&alias, "alias", "", "Alias to delete")
	return cmd
}
