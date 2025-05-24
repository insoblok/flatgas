package cmd

import "github.com/spf13/cobra"

func GetTxCommand() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction-related commands",
	}
	txCmd.AddCommand(GetTxStatusCommand())
	return txCmd
}
