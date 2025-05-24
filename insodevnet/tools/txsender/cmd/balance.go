package cmd

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
)

var alias string
var address string
var rpc string
var base string

var accountsBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check the ETH balance of an alias or address",
	RunE: func(cmd *cobra.Command, args []string) error {
		var resolved string
		var err error

		if alias != "" {
			resolved, err = internal.ResolveAddressOrAlias(base, alias)
			if err != nil {
				return fmt.Errorf("invalid alias: %w", err)
			}
		} else if address != "" {
			if !strings.HasPrefix(address, "0x") || len(address) != 42 {
				return fmt.Errorf("invalid address")
			}
			resolved = address
		} else {
			return fmt.Errorf("must provide either --alias or --address")
		}

		resolvedRPC := rpc
		if config, err := internal.LoadConfig(base); err == nil {
			if val, ok := config.RPCs[rpc]; ok {
				resolvedRPC = val
			}
		}
		client, err := ethclient.Dial(resolvedRPC)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC: %w", err)
		}
		defer client.Close()

		addr := common.HexToAddress(resolved)
		bal, err := client.BalanceAt(context.Background(), addr, nil)
		if err != nil {
			return fmt.Errorf("failed to get balance: %w", err)
		}

		eth := new(big.Float).Quo(new(big.Float).SetInt(bal), big.NewFloat(1e18))
		fmt.Printf("💰 Balance of %s: %s ETH\n", resolved, eth.Text('f', 6))
		return nil
	},
}

func init() {
	accountsBalanceCmd.Flags().StringVar(&alias, "alias", "", "Alias of account")
	accountsBalanceCmd.Flags().StringVar(&address, "address", "", "Raw Ethereum address")
	accountsBalanceCmd.Flags().StringVar(&rpc, "rpc", "http://localhost:8545", "RPC endpoint")
	accountsBalanceCmd.Flags().StringVar(&base, "base", ".", "Base path to flatgas root")
}

func GetAccountsBalanceCommand() *cobra.Command {
	return accountsBalanceCmd
}
