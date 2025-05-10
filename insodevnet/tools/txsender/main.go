package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"path/filepath"
)

type AccountMeta struct {
	Address  string
	Keyfile  string
	Password string
}

func getAccounts(base string) map[string]AccountMeta {

	var accounts = map[string]AccountMeta{
		"alice": {
			Address:  "0x3e7B248045C7B12B9eE5A484aF3eC7b4a3E37E87",
			Keyfile:  expand(base, "insodevnet/docker/single-validator/genesis/keys/devaccounts/UTC--2025-05-06T21-18-44.119194000Z--3e7b248045c7b12b9ee5a484af3ec7b4a3e37e87"),
			Password: "flatgas",
		},
		"bob": {
			Address:  "0xC7468BC7c7e0BA36fd369DF8B3eBf57ff5b0fC66",
			Keyfile:  expand(base, "insodevnet/docker/single-validator/genesis/keys/devaccounts/UTC--2025-05-06T21-18-47.394229000Z--c7468bc7c7e0ba36fd369df8b3ebf57ff5b0fc66"),
			Password: "flatgas",
		},
		"carol": {
			Address:  "0xC9ee165aA2D93C38527Cc72Ebb8ABabB9c830D36",
			Keyfile:  expand(base, "insodevnet/docker/single-validator/genesis/keys/devaccounts/UTC--2025-05-06T21-18-45.779238000Z--c9ee165aa2d93c38527cc72ebb8ababb9c830d36"),
			Password: "flatgas",
		},
	}
	return accounts
}

func expand(base, rel string) string {
	return filepath.Join(base, rel)
}

func main() {
	from := flag.String("from", "", "Alias of sender account (e.g. alice)")
	to := flag.String("to", "", "Recipient address")
	amount := flag.Float64("amount", 0, "Amount to send (in ETH)")
	rpc := flag.String("rpc", "http://localhost:8545", "RPC endpoint")
	base := flag.String("base", "", "Base path to Flatgas repo (default: current dir)")

	flag.Parse()
	resolvedBase := *base
	if resolvedBase == "" {
		resolvedBase, _ = os.Getwd()
	}

	if *from == "" || *to == "" || *amount <= 0 {
		fmt.Println("Usage: txsender --from <alias> --to <addr> --amount <eth> [--rpc <url>]")
		os.Exit(1)
	}

	fmt.Printf("🔄 Sending %f ETH from %s to %s via %s\n", *amount, *from, *to, *rpc)

	// ✅ Load account
	accounts := getAccounts(resolvedBase)
	account, ok := accounts[*from]
	if !ok {
		log.Fatalf("❌ Unknown alias: %s\n", *from)
	}

	keyjson, err := os.ReadFile(account.Keyfile)
	if err != nil {
		log.Fatalf("❌ Failed to read keyfile: %v\n", err)
	}

	key, err := keystore.DecryptKey(keyjson, account.Password)
	if err != nil {
		log.Fatalf("❌ Failed to decrypt key: %v\n", err)
	}

	fmt.Printf("🔓 Loaded account: %s (alias: %s)\n", key.Address.Hex(), *from)

	client, err := ethclient.Dial(*rpc)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RPC: %v\n", err)
	}
	defer client.Close()

	ctx := context.Background()

	nonce, err := client.PendingNonceAt(ctx, key.Address)
	if err != nil {
		log.Fatalf("❌ Failed to get nonce: %v\n", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to get chain ID: %v\n", err)
	}

	toAddr := common.HexToAddress(*to)
	amountWei := new(big.Int)
	amountWei.SetString(fmt.Sprintf("%.0f", *amount*1e18), 10) // ETH → wei

	gasLimit := uint64(21000) // standard tx
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to get gas price: %v\n", err)
	}

	tx := types.NewTransaction(nonce, toAddr, amountWei, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key.PrivateKey)
	if err != nil {
		log.Fatalf("❌ Failed to sign tx: %v\n", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("❌ Failed to send tx: %v\n", err)
	}

	fmt.Printf("📤 Transaction sent! Hash: %s\n", signedTx.Hash().Hex())
}
