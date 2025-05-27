package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/cmd/kvaccounts"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"time"
)

var kvAlias string
var kvPassword string

var kvaccountsCmd = &cobra.Command{
	Use:   "kvaccounts",
	Short: "Manage accounts using embedded KV store",
	Long:  "Provides versioned account creation, listing, rollback, and export backed by a local key-value store.",
}

var kvRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the most recent journaled action",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		dbPath := internal.GetDBFilePath(base)
		fmt.Println("🔧 Opening DB at:", dbPath)
		db, err := bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			return fmt.Errorf("failed to open db: %w", err)
		}
		defer db.Close()

		fmt.Println("🔍 Reading latest journal entry...")
		var lastKey []byte
		var lastEntry internal.JournalEntry
		err = db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("journal"))
			if bucket == nil {
				return fmt.Errorf("no journal entries to rollback")
			}
			c := bucket.Cursor()
			k, v := c.Last()
			if k == nil {
				return fmt.Errorf("journal is empty")
			}
			fmt.Printf("🧾 Selected rollback entry key: %s\n", k)
			lastKey = k
			return json.Unmarshal(v, &lastEntry)
		})
		if err != nil {
			return err
		}

		fmt.Printf("🔄 Entry to rollback: [%s] %s\n", lastEntry.Timestamp.Format("2006-01-02 15:04:05"), lastEntry.Alias)

		if lastEntry.Action != internal.ActionCreate {
			return fmt.Errorf("rollback for '%s' not implemented", lastEntry.Action)
		}

		fmt.Println("⚙️ Executing rollback...")
		err = db.Update(func(tx *bbolt.Tx) error {
			aliases := tx.Bucket([]byte("aliases"))
			if aliases == nil {
				return fmt.Errorf("aliases bucket not found")
			}
			if err := aliases.Delete([]byte(lastEntry.Alias)); err != nil {
				return err
			}

			journal := tx.Bucket([]byte("journal"))
			if journal == nil {
				return fmt.Errorf("journal bucket not found")
			}
			if err := journal.Delete(lastKey); err != nil {
				return err
			}

			rollbackEntry := internal.JournalEntry{
				Action:    internal.ActionRollback,
				Alias:     lastEntry.Alias,
				Timestamp: time.Now(),
				Data:      lastEntry.Data,
			}
			audit := tx.Bucket([]byte("auditlog"))
			if audit == nil {
				var err error
				audit, err = tx.CreateBucket([]byte("auditlog"))
				if err != nil {
					return fmt.Errorf("create auditlog bucket: %w", err)
				}
			}
			key := []byte(rollbackEntry.Timestamp.Format(time.RFC3339Nano))
			data, err := json.Marshal(rollbackEntry)
			if err != nil {
				return fmt.Errorf("marshal rollback entry: %w", err)
			}
			return audit.Put(key, data)
		})
		if err != nil {
			return err
		}
		fmt.Printf("↩️ Rolled back '%s'\n", lastEntry.Alias)
		return nil
	},
}

func GetKVAccountsCommand() *cobra.Command {
	kvaccountsCmd.PersistentFlags().String("base", ".", "Base path to flatgas repo")
	kvaccountsCmd.AddCommand(kvaccounts.GetCreateCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetListCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetHistoryCmd())
	kvaccountsCmd.AddCommand(kvRollbackCmd)
	kvaccountsCmd.AddCommand(kvaccounts.GetAuditCmd())
	return kvaccountsCmd
}
