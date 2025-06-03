package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Smart contract related operations",
}

var contractDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Compile and deploy a smart contract",
	Run: func(cmd *cobra.Command, args []string) {
		src, _ := cmd.Flags().GetString("src")

		fmt.Printf("📦 Compiling contract from %s...\n", src)

		// Check file exists
		if _, err := os.Stat(src); os.IsNotExist(err) {
			fmt.Printf("❌ File not found: %s\n", src)
			return
		}

		// Call solc to compile
		out, err := exec.Command("solc", "--combined-json", "abi,bin", src).Output()
		if err != nil {
			fmt.Printf("❌ Compilation failed: %s\n", err)
			return
		}

		var solcOutput struct {
			Contracts map[string]struct {
				ABI json.RawMessage `json:"abi"`
				Bin string          `json:"bin"`
			} `json:"contracts"`
		}

		if err := json.Unmarshal(out, &solcOutput); err != nil {
			fmt.Printf("❌ Failed to parse solc output: %s\n", err)
			return
		}

		if len(solcOutput.Contracts) == 0 {
			fmt.Println("⚠️ No contracts found in source.")
			return
		}

		for name, c := range solcOutput.Contracts {
			fmt.Printf("✅ Found contract: %s\n", name)
			fmt.Printf("📜 ABI: %s\n", string(c.ABI))
			fmt.Printf("🔢 Bytecode: %s... (%d bytes)\n", c.Bin[:32], len(c.Bin)/2)
		}
	},
}

func GetContractCommand() *cobra.Command {
	contractDeployCmd.Flags().String("src", "", "Path to the Solidity contract source file")
	contractDeployCmd.MarkFlagRequired("src")
	contractCmd.AddCommand(contractDeployCmd)
	return contractCmd
}
