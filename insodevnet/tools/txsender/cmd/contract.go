package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DeploymentResult struct {
	Contract   string    `json:"contract"`
	Source     string    `json:"source"`
	Network    string    `json:"network"`
	Address    string    `json:"address"`
	TxHash     string    `json:"txHash"`
	DeployedAt time.Time `json:"deployedAt"`
}

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Smart contract related operations",
}

var contractDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Compile and deploy a smart contract",
	Run: func(cmd *cobra.Command, args []string) {
		base, _ := cmd.Flags().GetString("base")
		dbPath := internal.GetAccountsDBFilePath(base)
		db, err := bbolt.Open(dbPath, 0600, nil)
		if err != nil {
			fmt.Errorf("failed to open DB: %w", err)
			return
		}
		defer db.Close()

		src, _ := cmd.Flags().GetString("src")
		outDir, _ := cmd.Flags().GetString("out")
		alias, _ := cmd.Flags().GetString("from")
		rpcURL, _ := cmd.Flags().GetString("rpc")
		password, _ := cmd.Flags().GetString("password")

		fmt.Printf("📦 Compiling contract from %s ...\n", src)
		cmdOut, err := exec.Command("solc", "--evm-version", "london", "--combined-json", "abi,bin", src).Output()
		if err != nil {
			fmt.Printf("❌ Compilation failed: %v\n", err)
			return
		}

		var solcOut struct {
			Contracts map[string]struct {
				ABI json.RawMessage `json:"abi"`
				Bin string          `json:"bin"`
			} `json:"contracts"`
		}

		if err := json.Unmarshal(cmdOut, &solcOut); err != nil {
			fmt.Printf("❌ Failed to parse solc output: %v\n", err)
			return
		}

		for name, contract := range solcOut.Contracts {
			fmt.Printf("✅ Found contract: %s\n", name)
			fmt.Printf("📜 ABI: %s\n", contract.ABI)
			fmt.Printf("🔢 Bytecode: %.20s... (%d bytes)\n", contract.Bin, len(contract.Bin)/2)

			bytecode := common.FromHex(contract.Bin)

			var record internal.AliasRecord

			err = db.View(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte("aliases"))
				if bucket == nil {
					return fmt.Errorf("aliases bucket not found")
				}

				data := bucket.Get([]byte(alias))
				if data == nil {
					return fmt.Errorf("alias not found: %s", alias)
				}
				return json.Unmarshal(data, &record)
			})
			if err != nil {
				fmt.Printf("❌ failed to read alias: %v\n", err)
				os.Exit(-1)
			}

			keyJSON, err := json.Marshal(record.Keystore)
			if err != nil {
				fmt.Printf("❌ failed to marshal keystore: %v\n", err)
				return
			}

			account, err := keystore.DecryptKey(keyJSON, password)
			if err != nil {
				fmt.Printf("❌ failed to decrypt key: %v\n", err)
				return
			}

			client, err := ethclient.Dial(rpcURL)
			if err != nil {
				log.Fatal(err)
			}
			defer client.Close()

			fromAddr := account.Address
			nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
			if err != nil {
				log.Fatal(err)
			}

			gasLimit := uint64(3000000)
			gasPrice := big.NewInt(1e9)
			value := big.NewInt(0)

			chainID, _ := client.NetworkID(context.Background())

			tx := types.NewContractCreation(nonce, value, gasLimit, gasPrice, bytecode)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), account.PrivateKey)
			if err != nil {
				log.Fatal(err)
			}

			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				log.Fatal(err)
			}

			txHash := signedTx.Hash().Hex()
			fmt.Printf("🚀 Deployment transaction sent. Hash: %s\n", txHash)

			receipt, err := waitForReceipt(client, txHash)
			if err != nil {
				log.Fatal(err)
			}

			addr := receipt.ContractAddress.Hex()
			fmt.Printf("✅ Contract deployed at: %s\n", addr)

			if outDir != "" {
				baseName := filepath.Base(name)
				baseName = strings.Split(baseName, ":")[1]
				targetDir := filepath.Join(outDir, baseName)
				os.MkdirAll(targetDir, 0755)

				os.WriteFile(filepath.Join(targetDir, baseName+".abi"), []byte(contract.ABI), 0644)
				os.WriteFile(filepath.Join(targetDir, baseName+".bin"), []byte(contract.Bin), 0644)

				meta := DeploymentResult{
					Contract:   baseName,
					Source:     src,
					Network:    "devnet",
					Address:    addr,
					TxHash:     txHash,
					DeployedAt: time.Now().UTC(),
				}
				data, _ := json.MarshalIndent(meta, "", "  ")
				os.WriteFile(filepath.Join(targetDir, baseName+".deploy.json"), data, 0644)
				fmt.Printf("💾 Written output to: %s\n", targetDir)
			}
		}
	},
}

func waitForReceipt(client *ethclient.Client, txHash string) (*types.Receipt, error) {
	ctx := context.Background()
	hash := common.HexToHash(txHash)
	for {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if err == nil {
			return receipt, nil
		}
		time.Sleep(1 * time.Second)
	}
}

func GetContractCommand() *cobra.Command {
	contractDeployCmd.Flags().String("src", "", "Path to the Solidity contract source file")
	contractDeployCmd.Flags().String("out", "", "Directory to output compiled contract artifacts")
	contractDeployCmd.Flags().String("from", "", "Alias name of sender account")
	contractDeployCmd.Flags().String("rpc", "", "Flatgas RPC endpoint (e.g., http://localhost:8545)")
	contractDeployCmd.Flags().String("password", "", "Password to decrypt key")
	contractDeployCmd.Flags().String("base", "", "Password to decrypt key")
	contractDeployCmd.MarkFlagRequired("rpc")
	contractDeployCmd.MarkFlagRequired("src")
	contractDeployCmd.MarkFlagRequired("from")
	contractDeployCmd.MarkFlagRequired("password")
	contractDeployCmd.MarkFlagRequired("base")
	contractCmd.AddCommand(contractDeployCmd)
	return contractCmd
}
