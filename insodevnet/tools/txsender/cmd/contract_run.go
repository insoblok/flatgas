package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var contractRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a method on a deployed smart contract",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 txsender contract run starting...")

		base, _ := cmd.Flags().GetString("base")
		contractDir, _ := cmd.Flags().GetString("dir")
		methodName, _ := cmd.Flags().GetString("method")
		methodArgs, _ := cmd.Flags().GetString("args")
		from, _ := cmd.Flags().GetString("from")
		rpcURL, _ := cmd.Flags().GetString("rpc")
		password, _ := cmd.Flags().GetString("password")
		gasLimit, _ := cmd.Flags().GetUint64("gas")

		info, err := os.Stat(contractDir)
		if err != nil || !info.IsDir() {
			fmt.Printf("❌ Invalid contract directory: %s\n", contractDir)
			os.Exit(1)
		}

		contractName := filepath.Base(filepath.Dir(contractDir))
		fmt.Printf("📦 Using contract: %s\n", contractName)

		abiFile := filepath.Join(contractDir, contractName+".abi")
		abiData, err := ioutil.ReadFile(abiFile)
		PrintIfErrorAndExit("Failed to read ABI file", err)

		contractABI, err := abi.JSON(strings.NewReader(string(abiData)))
		PrintIfErrorAndExit("Failed to read ABI ", err)

		fmt.Println("📜 Available methods:")
		for name, method := range contractABI.Methods {
			mut := "view"
			if !method.IsConstant() {
				mut = "nonpayable"
			}
			fmt.Printf(" - %s(%s) -> %s\n", name, method.Inputs, mut)
		}
		method, exists := contractABI.Methods[methodName]

		if !exists {
			fmt.Printf("❌ Method '%s' not found in contract ABI.\n", methodName)
			os.Exit(1)
		}

		var inputData []byte
		if len(method.Inputs) > 0 {
			if methodArgs == "" {
				fmt.Printf("❌ Method '%s' expects %d argument(s), but none were provided (--args).\n", methodName, len(method.Inputs))
				os.Exit(1)
			}
			var parsedArgs []interface{}
			err := json.Unmarshal([]byte(methodArgs), &parsedArgs)
			if err != nil {
				fmt.Printf("❌ Failed to parse --args as JSON: %v\n", err)
				os.Exit(1)
			}
			if len(parsedArgs) != len(method.Inputs) {
				fmt.Printf("❌ Method '%s' expects %d argument(s), but got %d.\n", methodName, len(method.Inputs), len(parsedArgs))
				os.Exit(1)
			}
			inputData, err = method.Inputs.Pack(parsedArgs...)
			PrintIfErrorAndExit("Failed to pack arguments", err)
		} else {
			if methodArgs != "" {
				fmt.Printf("⚠️ Method '%s' does not take arguments. Provided --args will be ignored.\n", methodName)
			}
			inputData, err = method.Inputs.Pack()
			PrintIfErrorAndExit("Failed to pack", err)
		}

		println(inputData)
		isView := method.StateMutability == "view" || method.StateMutability == "pure"
		isTransacted := !isView
		if isView {
			fmt.Printf("ℹ️  Method '%s' is read-only (no gas needed).\n", methodName)
		} else {
			fmt.Printf("⛽  Method '%s' is transacted (requires gas).\n", methodName)
		}

		client, err := ethclient.Dial(rpcURL)
		PrintIfErrorAndExit("Failed to connect to Ethereum RPC", err)
		defer client.Close()

		if isTransacted {
			if from == "" {
				fmt.Println("❌ This method modifies state and requires a sender (--from).")
				os.Exit(1)
			}

			if password == "" {
				fmt.Println("❌ --password is required.")
				os.Exit(1)
			}

			dbPath := internal.GetAccountsDBFilePath(base)
			db, err := bbolt.Open(dbPath, 0600, nil)
			PrintIfErrorAndExit("failed to open DB", err)
			defer db.Close()
			var record internal.AliasRecord
			err = db.View(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte("aliases"))
				if bucket == nil {
					return fmt.Errorf("Aliases bucket not found")
				}

				fmt.Println("Looking for alias", from)
				data := bucket.Get([]byte(from))
				if data == nil {
					return fmt.Errorf("Alias not found: %s", from)
				}
				return json.Unmarshal(data, &record)
			})
			PrintIfErrorAndExit("Failed to read alias", err)

			keyJSON, err := json.Marshal(record.Keystore)
			PrintIfErrorAndExit("Failed to marshal keystore", err)

			account, err := keystore.DecryptKey(keyJSON, password)
			PrintIfErrorAndExit("Failed to decrypt key", err)
			fmt.Println(account.Address.Hex())
		} else {

		}

		fmt.Println("📦 Parsed inputs:")
		fmt.Println("ContractDir:", contractDir)
		fmt.Println("Method:", methodName)
		fmt.Println("Args:", methodArgs)
		fmt.Println("From:", from)
		fmt.Println("RPC:", rpcURL)
		fmt.Println("Password:", password)
		fmt.Println("Gas:", gasLimit)
	},
}

func init() {
	contractRunCmd.Flags().String("base", ".", "Contract dir (required)")
	contractRunCmd.Flags().String("dir", "", "Contract dir (required)")
	contractRunCmd.Flags().String("method", "", "Method name to call (required)")
	contractRunCmd.Flags().String("args", "", "Constructor or method arguments as JSON array")
	contractRunCmd.Flags().String("from", "", "Alias of sender account")
	contractRunCmd.Flags().String("rpc", "", "RPC endpoint to connect to")
	contractRunCmd.Flags().String("password", "", "Password for account unlock")
	contractRunCmd.Flags().Uint64("gas", 3000000, "Optional gas limit")

	_ = contractRunCmd.MarkFlagRequired("dir")
	_ = contractRunCmd.MarkFlagRequired("method")
	_ = contractRunCmd.MarkFlagRequired("rpc")

	contractCmd.AddCommand(contractRunCmd)
}
