package cmd

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
)

var (
	fundFrom     string
	fundTo       string
	fundAmount   float64
	fundPassword string
	fundRPC      string
	fundSend     bool
)

var fundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Send ETH from a known account to another alias or address",
	RunE: func(cmd *cobra.Command, args []string) error {
		base, _ := cmd.Flags().GetString("base")
		rpcURL := fundRPC

		cfgData, err := internal.LoadConfig(base)
		if err == nil {
			if resolved, ok := cfgData.RPCs[fundRPC]; ok {
				rpcURL = resolved
			}
		}

		fromAddr, err := internal.ResolveAddressOrAlias(base, fundFrom)
		if err != nil {
			return fmt.Errorf("invalid --from: %w", err)
		}
		toAddr, err := internal.ResolveAddressOrAlias(base, fundTo)
		if err != nil {
			return fmt.Errorf("invalid --to: %w", err)
		}

		ks := keystore.NewKeyStore(filepath.Join(base, "wallet", "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
		var fromAccount accounts.Account
		found := false
		for _, acc := range ks.Accounts() {
			if acc.Address.Hex() == fromAddr {
				fromAccount = acc
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("account for address %s not found in keystore", fromAddr)
		}

		if err := ks.Unlock(fromAccount, fundPassword); err != nil {
			return fmt.Errorf("failed to unlock sender account: %w", err)
		}

		client, err := ethclient.Dial(rpcURL)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC: %w", err)
		}
		defer client.Close()

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get chain ID: %w", err)
		}

		nonce, err := client.PendingNonceAt(context.Background(), fromAccount.Address)
		if err != nil {
			return fmt.Errorf("failed to get nonce: %w", err)
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get gas price: %w", err)
		}

		value := big.NewInt(int64(fundAmount * 1e18))
		to := common.HexToAddress(toAddr)
		var data []byte
		gasLimit := uint64(21000)
		tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

		signedTx, err := ks.SignTx(fromAccount, tx, chainID)
		if err != nil {
			return fmt.Errorf("failed to sign tx: %w", err)
		}

		if !fundSend {
			fmt.Println("🧪 Dry run successful. Transaction prepared but not broadcast.")
			fmt.Printf("🔗 Would send %f ETH from %s to %s\n", fundAmount, fromAddr, toAddr)
			fmt.Printf("📦 Tx hash (preview): %s\n", signedTx.Hash().Hex())
			return nil
		}

		if err := client.SendTransaction(context.Background(), signedTx); err != nil {
			return fmt.Errorf("failed to send tx: %w", err)
		}

		fmt.Printf("📤 Sent %f ETH from %s to %s\n", fundAmount, fromAddr, toAddr)
		fmt.Printf("🔗 Tx hash: %s\n", signedTx.Hash().Hex())
		return nil
	},
}

func init() {
	fundCmd.Flags().StringVar(&fundFrom, "from", "", "Alias or address of sender")
	fundCmd.Flags().StringVar(&fundTo, "to", "", "Alias or address of recipient")
	fundCmd.Flags().Float64Var(&fundAmount, "amount", 0, "Amount in ETH to send")
	fundCmd.Flags().StringVar(&fundPassword, "password", "", "Password for sender account")
	fundCmd.Flags().StringVar(&fundRPC, "rpc", "http://localhost:8545", "RPC endpoint URL or alias")
	fundCmd.Flags().BoolVar(&fundSend, "send", false, "Broadcast the transaction (default is dry run)")
	fundCmd.Flags().String("base", ".", "Base path to flatgas root")
}

func GetFundCommand() *cobra.Command {
	return fundCmd
}
