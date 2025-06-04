package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"
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

		contractDir, _ := cmd.Flags().GetString("dir")
		methodName, _ := cmd.Flags().GetString("method")
		argsJSON, _ := cmd.Flags().GetString("args")
		from, _ := cmd.Flags().GetString("from")
		rpc, _ := cmd.Flags().GetString("rpc")
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

		isView := method.StateMutability == "view" || method.StateMutability == "pure"

		if isView {
			fmt.Printf("ℹ️  Method '%s' is read-only (no gas needed).\n", methodName)
		} else {
			fmt.Printf("⛽  Method '%s' is transacted (requires gas).\n", methodName)
		}

		fmt.Println("📦 Parsed inputs:")
		fmt.Println("ContractDir:", contractDir)
		fmt.Println("Method:", methodName)
		fmt.Println("Args:", argsJSON)
		fmt.Println("From:", from)
		fmt.Println("RPC:", rpc)
		fmt.Println("Password:", password)
		fmt.Println("Gas:", gasLimit)
	},
}

func init() {
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
