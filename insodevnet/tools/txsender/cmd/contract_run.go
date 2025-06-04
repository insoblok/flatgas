package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/insoblok/flatgas/insodevnet/tools/txsender/internal"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"io/ioutil"
	"log"
	"math/big"
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

		fmt.Printf("🔍 Base path: %s\n", base)
		fmt.Printf("📂 Contract directory: %s\n", contractDir)

		info, err := os.Stat(contractDir)
		if err != nil || !info.IsDir() {
			fmt.Printf("❌ Invalid contract directory: %s\n", contractDir)
			os.Exit(1)
		}

		contractName := filepath.Base(filepath.Dir(contractDir))
		fmt.Printf("📦 Using contract: %s\n", contractName)

		metaPath := filepath.Join(contractDir, contractName+".deploy.json")
		fmt.Printf("📄 Reading metadata from: %s\n", metaPath)
		metaData, err := os.ReadFile(metaPath)
		PrintIfErrorAndExit("❌ Failed to read metadata JSON", err)

		var meta struct {
			Address string `json:"address"`
		}
		err = json.Unmarshal(metaData, &meta)
		PrintIfErrorAndExit("❌ Failed to parse metadata JSON", err)
		contractAddress := common.HexToAddress(meta.Address)
		fmt.Printf("🏠 Contract address: %s\n", contractAddress.Hex())

		abiFile := filepath.Join(contractDir, contractName+".abi")
		fmt.Printf("📄 Reading ABI from: %s\n", abiFile)
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

		fmt.Printf("🔧 Selected method: %s\n", method.Name)
		fmt.Printf("🔧 Method mutability: %s\n", method.StateMutability)

		var inputData []byte
		if len(method.Inputs) > 0 {
			fmt.Printf("🧩 Method requires %d argument(s).\n", len(method.Inputs))
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
			fmt.Printf("🧩 Parsed arguments: %v\n", parsedArgs)
			if len(parsedArgs) != len(method.Inputs) {
				fmt.Printf("❌ Method '%s' expects %d argument(s), but got %d.\n", methodName, len(method.Inputs), len(parsedArgs))
				os.Exit(1)
			}
			fmt.Printf("🧪 Raw parsedArgs type check:\n")
			for i, arg := range parsedArgs {
				fmt.Printf("  [%d] Type: %T, Value: %#v\n", i, arg, arg)
			}

			var typedArgs []interface{}
			for i, arg := range parsedArgs {
				expected := method.Inputs[i].Type.String()
				switch expected {
				case "address":
					addrStr, ok := arg.(string)
					if !ok {
						log.Fatalf("Argument %d: expected string for address", i)
					}
					typedArgs = append(typedArgs, common.HexToAddress(addrStr))

				case "uint256":
					// You could also use json.Number + big.Int here for safer parsing
					floatVal, ok := arg.(float64)
					if !ok {
						log.Fatalf("Argument %d: expected number for uint256", i)
					}
					typedArgs = append(typedArgs, big.NewInt(int64(floatVal)))

				case "bool":
					boolVal, ok := arg.(bool)
					if !ok {
						log.Fatalf("Argument %d: expected bool", i)
					}
					typedArgs = append(typedArgs, boolVal)

				case "string":
					strVal, ok := arg.(string)
					if !ok {
						log.Fatalf("Argument %d: expected string", i)
					}
					typedArgs = append(typedArgs, strVal)

				default:
					log.Fatalf("Unsupported argument type: %s", expected)
				}
			}
			argData, err := method.Inputs.Pack(typedArgs...)
			PrintIfErrorAndExit("Failed to pack arguments", err)
			inputData = append(method.ID, argData...)
		} else {
			if methodArgs != "" {
				fmt.Printf("⚠️ Method '%s' does not take arguments. Provided --args will be ignored.\n", methodName)
			}
			argData, err := method.Inputs.Pack()
			inputData = append(method.ID, argData...)

			PrintIfErrorAndExit("Failed to pack", err)
		}

		fmt.Printf("📤 Packed input data: 0x%x\n", inputData)

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
			fmt.Printf("📁 Account DB path: %s\n", dbPath)
			db, err := bbolt.Open(dbPath, 0600, nil)
			PrintIfErrorAndExit("failed to open DB", err)
			defer db.Close()
			var record internal.AliasRecord
			err = db.View(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket([]byte("aliases"))
				if bucket == nil {
					return fmt.Errorf("Aliases bucket not found")
				}
				fmt.Println("🔍 Looking for alias", from)
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
			fmt.Println("🔓 Account address:", account.Address.Hex())
			// Build and send the transaction
			ctx := context.Background()

			nonce, err := client.PendingNonceAt(ctx, account.Address)
			PrintIfErrorAndExit("Failed to get account nonce", err)

			chainID, err := client.NetworkID(ctx)
			PrintIfErrorAndExit("Failed to get chain ID", err)

			gasPrice, err := client.SuggestGasPrice(ctx)
			PrintIfErrorAndExit("Failed to get gas price", err)

			tx := types.NewTransaction(
				nonce,
				contractAddress,
				big.NewInt(0), // value = 0 for method call
				gasLimit,
				gasPrice,
				inputData,
			)

			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), account.PrivateKey)
			PrintIfErrorAndExit("Failed to sign transaction", err)

			err = client.SendTransaction(ctx, signedTx)
			PrintIfErrorAndExit("Failed to send transaction", err)

			fmt.Printf("✅ Transaction sent: %s\n", signedTx.Hash().Hex())

		} else {
			msg := ethereum.CallMsg{
				To:   &contractAddress,
				Data: inputData,
			}

			ctx := context.Background()
			result, err := client.CallContract(ctx, msg, nil)
			if err != nil {
				log.Fatalf("CallContract error: %v", err)
			}

			fmt.Printf("📥 Raw return data: 0x%x\n", result)
			outputs := method.Outputs
			values, err := outputs.Unpack(result)
			if err != nil {
				fmt.Println("Failed to unpack result: %v", err)
			}

			fmt.Printf("✅ Method call result: %v\n", values)
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
