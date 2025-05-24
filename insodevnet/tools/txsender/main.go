package main

import (
	"fmt"
	"log"
	"os"

	"github.com/insoblok/flatgas/insodevnet/tools/txsender/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "txsender",
		Short: "Flatgas devnet CLI for account and transaction management",
	}

	rootCmd.AddCommand(
		cmd.GetConfigCommand(),
		cmd.GetAccountsCommand(),
		cmd.GetFundCommand(),
		cmd.GetTxCommand(),
		cmd.GetNodeCommand(),
	)
	fmt.Println("🚀 txsender CLI starting...")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("❌ Error: %v", err)
		os.Exit(1)
	}
}
