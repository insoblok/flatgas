package cmd

import (
	"fmt"
	cfg "github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
)

var (
	rpcName        string
	rpcURL         string
	defaultRPCName string
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

var listRpcsCmd = &cobra.Command{
	Use:   "list-rpcs",
	Short: "List configured RPC aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		cfgData, err := cfg.LoadConfig(base)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Println("📡 Configured RPCs:")
		for name, url := range cfgData.RPCs {
			marker := ""
			if name == cfgData.DefaultRPC {
				marker = " ✅"
			}
			fmt.Printf("  %s => %s%s\n", name, url, marker)
		}
		return nil
	},
}

var setDefaultRpcCmd = &cobra.Command{
	Use:   "set-default-rpc",
	Short: "Set the default RPC to use for all commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		if defaultRPCName == "" {
			return fmt.Errorf("--name is required")
		}

		cfgData, err := cfg.LoadConfig(base)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if _, exists := cfgData.RPCs[defaultRPCName]; !exists {
			return fmt.Errorf("RPC name '%s' not found in config", defaultRPCName)
		}

		cfgData.DefaultRPC = defaultRPCName
		if err := cfg.SaveConfig(base, cfgData); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Default RPC set to '%s'\n", defaultRPCName)
		return nil
	},
}

func init() {
	addRpcCmd.Flags().StringVar(&rpcName, "name", "", "Alias for the RPC node")
	addRpcCmd.Flags().StringVar(&rpcURL, "url", "", "RPC endpoint URL")
	addRpcCmd.Flags().String("base", ".", "Base path to flatgas root")

	listRpcsCmd.Flags().String("base", ".", "Base path to flatgas root")

	setDefaultRpcCmd.Flags().StringVar(&defaultRPCName, "name", "", "Alias of RPC to set as default")
	setDefaultRpcCmd.Flags().String("base", ".", "Base path to flatgas root")

	configCmd.AddCommand(addRpcCmd)
	configCmd.AddCommand(listRpcsCmd)
	configCmd.AddCommand(setDefaultRpcCmd)
}

func GetConfigCommand() *cobra.Command {
	return configCmd
}
