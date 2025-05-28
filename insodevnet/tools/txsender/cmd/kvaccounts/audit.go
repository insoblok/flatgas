package kvaccounts

import (
	"encoding/json"
	"fmt"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"sort"
)

var kvAuditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Display all actions from the immutable auditlog",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		dbPath := internal.GetAccountsDBFilePath(base)
		db, err := bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			return fmt.Errorf("failed to open db: %w", err)
		}
		defer db.Close()

		var entries []internal.JournalEntry
		err = db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("auditlog"))
			if bucket == nil {
				fmt.Println("No audit log entries found.")
				return nil
			}
			return bucket.ForEach(func(k, v []byte) error {
				var entry internal.JournalEntry
				if err := json.Unmarshal(v, &entry); err != nil {
					return err
				}
				entries = append(entries, entry)
				return nil
			})
		})
		if err != nil {
			return err
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.After(entries[j].Timestamp)
		})
		fmt.Println("🧾 Audit Log:")
		for _, e := range entries {
			fmt.Printf("- [%s] %s %s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Action, e.Alias)
		}
		return nil
	},
}

func GetAuditCmd() *cobra.Command {
	return kvAuditCmd
}
