package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var txHash string
var txRPC string

var txStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of a transaction by hash",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !strings.HasPrefix(txHash, "0x") || len(txHash) != 66 {
			return fmt.Errorf("invalid transaction hash: %s", txHash)
		}

		client, err := ethclient.Dial(txRPC)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC: %w", err)
		}
		defer client.Close()

		hash := common.HexToHash(txHash)
		receipt, err := client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			return fmt.Errorf("failed to get receipt: %w", err)
		}

		fmt.Printf("✅ TX is mined\n")
		fmt.Printf("🔗 Hash: %s\n", hash.Hex())
		fmt.Printf("📦 Block: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("⛽ Gas Used: %d\n", receipt.GasUsed)
		fmt.Printf("🧾 Status: %d\n", receipt.Status)

		result, _ := json.MarshalIndent(receipt, "", "  ")
		log.Printf("\n📋 Full Receipt:\n%s\n", result)
		return nil
	},
}

func init() {
	txStatusCmd.Flags().StringVar(&txHash, "tx", "", "Transaction hash")
	txStatusCmd.Flags().StringVar(&txRPC, "rpc", "http://localhost:8545", "RPC endpoint")
	txStatusCmd.MarkFlagRequired("tx")
}

func GetTxStatusCommand() *cobra.Command {
	return txStatusCmd
}
