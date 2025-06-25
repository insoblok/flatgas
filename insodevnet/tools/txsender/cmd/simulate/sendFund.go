package simulate

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"path/filepath"

	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
)

func SendFunds(base, from, fromPassword string, to string, value *big.Int, rpcURL string) error {
	resolvedFrom := MustResolve(base, from)
	resolvedTo := MustResolve(base, to)
	log.Printf("Would send %s wei from %s to %s rpc %s\n", value.String(), resolvedFrom.Hex(), resolvedTo.Hex(), rpcURL)

	ks := keystore.NewKeyStore(filepath.Join(base, "wallet", "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
	var fromAccount accounts.Account
	found := false
	for _, acc := range ks.Accounts() {
		if acc.Address.Hex() == resolvedFrom.Hex() {
			fromAccount = acc
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("account for address %s not found in keystore", resolvedFrom)
	}

	if err := ks.Unlock(fromAccount, fromPassword); err != nil {
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

	var data []byte
	gasLimit := uint64(21000)
	tx := types.NewTransaction(nonce, resolvedTo, value, gasLimit, gasPrice, data)

	signedTx, err := ks.SignTx(fromAccount, tx, chainID)
	if err != nil {
		return fmt.Errorf("failed to sign tx: %w", err)
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	fmt.Printf("📤 Sent %f ETH from %s to %s\n", value, resolvedFrom, resolvedTo)
	fmt.Printf("🔗 Tx hash: %s\n", signedTx.Hash().Hex())
	return nil

}

func MustResolve(base, input string) common.Address {
	addr, err := internal.ResolveAddressOrAlias(base, input)
	if err != nil {
		log.Fatalf("❌ Could not resolve %s: %v", input, err)
	}
	return common.HexToAddress(addr)
}
