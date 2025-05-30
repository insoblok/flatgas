package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	cfg "github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"strconv"
)

var nodeRPC string

var nodeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display information about the connected inso node",
	RunE: func(cmd *cobra.Command, args []string) error {
		resolvedRPC := nodeRPC
		base, _ := cmd.Flags().GetString("base")

		fmt.Printf("what is base %s\n", base)
		if config, err := cfg.LoadConfig(base); err == nil {
			if val, ok := config.RPCs[nodeRPC]; ok {
				resolvedRPC = val
			}
		}
		client, err := ethclient.Dial(resolvedRPC)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC: %w", err)
		}
		defer client.Close()

		ctx := context.Background()

		blockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			return fmt.Errorf("failed to get block number: %w", err)
		}

		chainID, err := client.ChainID(ctx)
		if err != nil {
			return fmt.Errorf("failed to get chain ID: %w", err)
		}

		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return fmt.Errorf("failed to get gas price: %w", err)
		}

		peerCount := "N/A"
		if netClient := client.Client(); netClient != nil {
			var countHex string
			err = netClient.CallContext(ctx, &countHex, "net_peerCount")
			if err == nil {
				countInt, _ := strconv.ParseInt(countHex[2:], 16, 64)
				peerCount = fmt.Sprintf("%d", countInt)
			}
		}

		fmt.Println("🖥️  Node Info")
		fmt.Printf("🧬 Chain ID: %s\n", chainID.String())
		fmt.Printf("📦 Latest Block: %d\n", blockNumber)
		fmt.Printf("⛽ Gas Price: %s wei\n", gasPrice.String())
		fmt.Printf("🤝 Peers: %s\n", peerCount)
		return nil
	},
}

func init() {
	nodeInfoCmd.Flags().String("base", ".", "Base path to flatgas root")
	nodeInfoCmd.Flags().StringVar(&nodeRPC, "rpc", "http://localhost:8545", "RPC endpoint")
}

func GetNodeInfoCommand() *cobra.Command {
	return nodeInfoCmd
}
