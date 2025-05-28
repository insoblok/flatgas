package kvaccounts

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var kvAlias string
var kvPassword string

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
			Metadata: map[string]string{"tag": "created"},
			Created:  time.Now(),
			Updated:  time.Now(),
		}

		dbPath := internal.GetAccountsDBFilePath(base)
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

		fmt.Println("📝 Preparing create audit log entry...")
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

func GetCreateCmd() *cobra.Command {
	kvCreateCmd.Flags().StringVar(&kvAlias, "alias", "", "Alias for the new account")
	kvCreateCmd.Flags().StringVar(&kvPassword, "password", "", "Password to encrypt key")
	return kvCreateCmd
}
