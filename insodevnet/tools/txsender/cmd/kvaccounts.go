package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var kvaccountsCmd = &cobra.Command{
	Use:   "kvaccounts",
	Short: "Manage accounts using embedded KV store",
	Long:  "Provides versioned account creation, listing, rollback, and export backed by a local key-value store.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧪 kvaccounts command is wired correctly!")
	},
}

func GetKVAccountsCommand() *cobra.Command {
	kvaccountsCmd.PersistentFlags().String("base", ".", "Base path to flatgas repo")
	return kvaccountsCmd
}
