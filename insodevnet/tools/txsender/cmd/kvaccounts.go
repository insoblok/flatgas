package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/cmd/kvaccounts"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var kvAlias string
var kvPassword string

var kvaccountsCmd = &cobra.Command{
	Use:   "kvaccounts",
	Short: "Manage accounts using embedded KV store",
	Long:  "Provides versioned account creation, listing, rollback, and export backed by a local key-value store.",
}

var kvCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new alias + encrypted account in kvstore",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		if kvAlias == "" || kvPassword == "" {
			return fmt.Errorf("--alias and --password are required")
		}

		tempKeystore := filepath.Join(os.TempDir(), "keystore-tmp")
		os.MkdirAll(tempKeystore, 0700)
		defer os.RemoveAll(tempKeystore)

		ks := keystore.NewKeyStore(tempKeystore, keystore.StandardScryptN, keystore.StandardScryptP)

		acct, err := ks.NewAccount(kvPassword)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		keyData, err := os.ReadFile(acct.URL.Path)
		if err != nil {
			return fmt.Errorf("failed to read keystore file: %w", err)
		}
		var keyRaw map[string]interface{}
		if err := json.Unmarshal(keyData, &keyRaw); err != nil {
			return fmt.Errorf("invalid keystore JSON: %w", err)
		}

		record := internal.AliasRecord{
			Alias:    internal.Alias(kvAlias),
			Address:  internal.KvAddress(acct.Address.Hex()),
			Keystore: keyRaw,
			Metadata: map[string]interface{}{"tags": []string{"created"}},
			Created:  time.Now(),
			Updated:  time.Now(),
		}

		dbPath := internal.GetDBFilePath(base)
		db, err := bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			return fmt.Errorf("failed to open db: %w", err)
		}
		defer db.Close()

		if err := db.Update(func(tx *bbolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("aliases"))
			if err != nil {
				return err
			}
			data, err := json.Marshal(record)
			if err != nil {
				return err
			}
			return bucket.Put([]byte(kvAlias), data)
		}); err != nil {
			return fmt.Errorf("failed to write alias to db: %w", err)
		}

		entry := internal.JournalEntry{
			Action:    internal.ActionCreate,
			Alias:     record.Alias,
			Timestamp: time.Now(),
			Data:      &record,
		}
		if err := internal.WriteJournalEntry(db, entry); err != nil {
			log.Fatalf("failed to write journal entry: %v", err)
		}
		if err := internal.WriteAuditLogEntry(db, entry); err != nil {
			log.Fatalf("failed to write audit log entry: %v", err)
		}

		fmt.Printf("✅ Created alias '%s' → %s\n", record.Alias, record.Address)
		return nil
	},
}

var kvHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View recent journal entries from the kvstore",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		dbPath := internal.GetDBFilePath(base)
		db, err := bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			return fmt.Errorf("failed to open db: %w", err)
		}
		defer db.Close()

		var entries []internal.JournalEntry
		err = db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("journal"))
			if bucket == nil {
				fmt.Println("No journal entries found.")
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
		fmt.Println("🕓 Journal History:")
		for _, e := range entries {
			fmt.Printf("- [%s] %s %s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Action, e.Alias)
		}
		return nil
	},
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

var kvAuditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Display all actions from the immutable auditlog",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		dbPath := internal.GetDBFilePath(base)
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

func GetKVAccountsCommand() *cobra.Command {
	kvaccountsCmd.PersistentFlags().String("base", ".", "Base path to flatgas repo")
	kvCreateCmd.Flags().StringVar(&kvAlias, "alias", "", "Alias for the new account")
	kvCreateCmd.Flags().StringVar(&kvPassword, "password", "", "Password to encrypt key")
	kvaccountsCmd.AddCommand(kvCreateCmd)
	kvaccountsCmd.AddCommand(kvaccounts.GetListCmd())
	kvaccountsCmd.AddCommand(kvHistoryCmd)
	kvaccountsCmd.AddCommand(kvRollbackCmd)
	kvaccountsCmd.AddCommand(kvAuditCmd)
	return kvaccountsCmd
}
