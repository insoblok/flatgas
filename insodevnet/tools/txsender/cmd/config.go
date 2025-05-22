package cmd

import (
	"fmt"
	cfg "github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
)

var (
	rpcName string
	rpcURL  string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage txsender configuration",
}

var addRpcCmd = &cobra.Command{
	Use:   "add-rpc",
	Short: "Add or update an RPC alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		if rpcName == "" || rpcURL == "" {
			return fmt.Errorf("both --name and --url are required")
		}

		cfgData, err := cfg.LoadConfig(base)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfgData.RPCs[rpcName] = rpcURL

		if err := cfg.SaveConfig(base, cfgData); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ RPC '%s' set to %s\n", rpcName, rpcURL)
		return nil
	},
}

func init() {
	addRpcCmd.Flags().StringVar(&rpcName, "name", "", "Alias for the RPC node")
	addRpcCmd.Flags().StringVar(&rpcURL, "url", "", "RPC endpoint URL")
	addRpcCmd.Flags().String("base", ".", "Base path to flatgas root")

	configCmd.AddCommand(addRpcCmd)
}

func GetConfigCommand() *cobra.Command {
	return configCmd
}
