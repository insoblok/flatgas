package kvaccounts

import (
	"fmt"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"time"
)

var kvMetaKey string
var kvMetaValue string

func GetMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Manage metadata on kvaccounts",
	}

	cmd.AddCommand(getMetaAddCmd())
	cmd.AddCommand(getMetaDeleteCmd())
	cmd.AddCommand(getMetaListCmd())

	return cmd
}

func getMetaAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add or update metadata key-value",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			if kvAlias == "" || kvMetaKey == "" || kvMetaValue == "" {
				return fmt.Errorf("alias, key, and value must be provided")
			}
			dbPath := internal.GetAccountsDBFilePath(base)
			fmt.Printf("📁 DB: %s\n", dbPath)
			return updateMetadata(dbPath, kvAlias, kvMetaKey, kvMetaValue, "add")
		},
	}
	cmd.Flags().StringVar(&kvMetaKey, "key", "", "Metadata key")
	cmd.Flags().StringVar(&kvMetaValue, "value", "", "Metadata value")
	cmd.Flags().StringVar(&kvAlias, "alias", "", "Alias of the account")
	return cmd
}

func getMetaDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a metadata key",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			if kvAlias == "" || kvMetaKey == "" {
				return fmt.Errorf("alias and key must be provided")
			}
			dbPath := internal.GetAccountsDBFilePath(base)
			return updateMetadata(dbPath, kvAlias, kvMetaKey, "", "delete")
		},
	}
	cmd.Flags().StringVar(&kvMetaKey, "key", "", "Metadata key")
	cmd.Flags().StringVar(&kvAlias, "alias", "", "Alias of the account")
	return cmd
}

func getMetaListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List metadata for alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			if kvAlias == "" {
				return fmt.Errorf("alias must be provided")
			}
			dbPath := internal.GetAccountsDBFilePath(base)
			db, err := bbolt.Open(dbPath, 0600, nil)
			if err != nil {
				return fmt.Errorf("failed to open db: %w", err)
			}
			defer db.Close()

			record, err := internal.ReadAlias(db, kvAlias)
			if err != nil {
				return err
			}
			fmt.Printf("🔎 Metadata for %s:\n", kvAlias)
			for k, v := range record.Metadata {
				fmt.Printf("  %s: %s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&kvAlias, "alias", "", "Alias of the account")
	return cmd
}

func updateMetadata(dbPath, alias, key, value, op string) error {
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()
	return internal.WithUpdateAlias(db, alias, func(record *internal.AliasRecord) error {
		if record.Metadata == nil {
			record.Metadata = map[string]string{}
		}
		switch op {
		case "add":
			record.Metadata[key] = value
			fmt.Printf("📝 Set %s = %s\n", key, value)
		case "delete":
			if _, ok := record.Metadata[key]; !ok {
				return fmt.Errorf("key not found: %s", key)
			}
			delete(record.Metadata, key)
			fmt.Printf("❌ Deleted key: %s\n", key)
		}
		record.Updated = time.Now()
		return nil
	})
}
