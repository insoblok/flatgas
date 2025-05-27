package cmd

import (
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/cmd/kvaccounts"
	"github.com/spf13/cobra"
)

var kvAlias string
var kvPassword string

var kvaccountsCmd = &cobra.Command{
	Use:   "kvaccounts",
	Short: "Manage accounts using embedded KV store",
	Long:  "Provides versioned account creation, listing, rollback, and export backed by a local key-value store.",
}

func GetKVAccountsCommand() *cobra.Command {
	kvaccountsCmd.PersistentFlags().String("base", ".", "Base path to flatgas repo")
	kvaccountsCmd.AddCommand(kvaccounts.GetCreateCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetListCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetHistoryCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetRollbackCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetAuditCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetDeleteCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetExportCmd())
	kvaccountsCmd.AddCommand(kvaccounts.GetMetadataCmd())
	return kvaccountsCmd
}
