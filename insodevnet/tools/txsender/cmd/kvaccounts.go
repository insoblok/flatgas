package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
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

		fmt.Printf("✅ Created alias '%s' → %s\n", record.Alias, record.Address)
		return nil
	},
}

var kvListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases stored in the kvstore",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		dbPath := internal.GetDBFilePath(base)
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

func GetKVAccountsCommand() *cobra.Command {
	kvaccountsCmd.PersistentFlags().String("base", ".", "Base path to flatgas repo")
	kvCreateCmd.Flags().StringVar(&kvAlias, "alias", "", "Alias for the new account")
	kvCreateCmd.Flags().StringVar(&kvPassword, "password", "", "Password to encrypt key")
	kvaccountsCmd.AddCommand(kvCreateCmd)
	kvaccountsCmd.AddCommand(kvListCmd)
	return kvaccountsCmd
}
