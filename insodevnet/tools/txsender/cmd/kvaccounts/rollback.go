package kvaccounts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

func GetRollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback the latest change (create/delete/update-meta) to aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			dbPath := internal.GetDBFilePath(base)
			fmt.Printf("🔧 Opening DB at: %s\n", dbPath)

			db, err := bbolt.Open(dbPath, 0600, nil)
			if err != nil {
				return fmt.Errorf("failed to open db: %w", err)
			}
			defer db.Close()

			fmt.Println("🔍 Reading latest journal entry...")
			var lastKey []byte
			var lastEntry internal.JournalEntry

			err = db.View(func(tx *bbolt.Tx) error {
				journal := tx.Bucket([]byte("journal"))
				if journal == nil {
					return fmt.Errorf("journal bucket not found")
				}
				c := journal.Cursor()
				k, v := c.Last()
				if k == nil {
					return fmt.Errorf("no journal entries found")
				}
				lastKey = k
				return json.Unmarshal(v, &lastEntry)
			})
			if err != nil {
				return err
			}

			fmt.Printf("💾 Selected rollback entry key: %s\n", string(lastKey))
			fmt.Printf("🔄 Entry to rollback: [%s] %s (%s)\n", lastEntry.Timestamp.Format("2006-01-02 15:04:05"), lastEntry.Alias, lastEntry.Action)
			fmt.Println("⚙️ Executing rollback...")

			err = db.Update(func(tx *bbolt.Tx) error {
				aliases := tx.Bucket([]byte("aliases"))
				if aliases == nil {
					return fmt.Errorf("aliases bucket not found")
				}

				switch lastEntry.Action {
				case internal.ActionCreate:
					fmt.Printf("🗑️ Deleting alias '%s'...\n", lastEntry.Alias)
					if err := aliases.Delete([]byte(lastEntry.Alias)); err != nil {
						return fmt.Errorf("failed to delete alias during rollback: %w", err)
					}
				case internal.ActionDelete:
					fmt.Printf("♻️ Restoring alias '%s'...\n", lastEntry.Alias)
					fmt.Println("🔧 Marshalling alias record...")
					data, err := json.Marshal(lastEntry.Data)
					if err != nil {
						return fmt.Errorf("failed to marshal restored alias: %w", err)
					}
					fmt.Println("✅ Marshal complete. Proceeding to PUT...")
					fmt.Println("🔧 Writing alias record to DB...")
					if err := aliases.Put([]byte(lastEntry.Alias), data); err != nil {
						return fmt.Errorf("failed to restore alias '%s': %w", lastEntry.Alias, err)
					}
					fmt.Println("✅ Alias restored in DB.")
				case internal.ActionUpdateMeta:
					data, err := json.Marshal(lastEntry.Data)
					if err != nil {
						return fmt.Errorf("failed to marshal updated alias: %w", err)
					}
					if err := aliases.Put([]byte(lastEntry.Alias), data); err != nil {
						return fmt.Errorf("failed to restore alias during rollback: %w", err)
					}

				default:
					return fmt.Errorf("rollback not supported for action: %s", lastEntry.Action)
				}

				journal := tx.Bucket([]byte("journal"))
				if journal == nil {
					return fmt.Errorf("journal bucket not found")
				}
				if err := journal.Delete(lastKey); err != nil {
					return fmt.Errorf("failed to delete journal entry: %w", err)
				}

				fmt.Println("📝 Preparing rollback audit log entry...")
				rollbackEntry := internal.JournalEntry{
					Action:    internal.ActionRollback,
					Alias:     lastEntry.Alias,
					Timestamp: time.Now(),
					Data:      lastEntry.Data,
				}
				fmt.Println("📋 Writing audit log entry...")
				if err := internal.WriteTxAuditLogEntry(tx, rollbackEntry); err != nil {
					return fmt.Errorf("failed to write audit log: %w", err)
				}
				fmt.Println("✅ Audit log entry written.")
				return nil
			})
			if err != nil {
				return err
			}

			fmt.Printf("↩️ Rolled back '%s'\n", lastEntry.Alias)
			return nil
		},
	}
}
