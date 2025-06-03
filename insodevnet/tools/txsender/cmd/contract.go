package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
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
		src, _ := cmd.Flags().GetString("src")
		outDir, _ := cmd.Flags().GetString("out")

		fmt.Println("📦 Compiling contract from", src, "...")
		cmdOut, err := exec.Command("solc", "--combined-json", "abi,bin", src).Output()
		if err != nil {
			fmt.Println("❌ Compilation failed:", err)
			return
		}

		var solcOut struct {
			Contracts map[string]struct {
				ABI json.RawMessage `json:"abi"`
				Bin string          `json:"bin"`
			} `json:"contracts"`
		}

		if err := json.Unmarshal(cmdOut, &solcOut); err != nil {
			fmt.Println("❌ Failed to parse solc output:", err)
			return
		}

		for name, contract := range solcOut.Contracts {
			fmt.Println("✅ Found contract:", name)
			fmt.Println("📜 ABI:", contract.ABI)
			fmt.Printf("🔢 Bytecode: %.20s... (%d bytes)\n", contract.Bin, len(contract.Bin)/2)

			if outDir != "" {
				baseName := filepath.Base(name)
				baseName = strings.Split(baseName, ":")[1] // From contracts/Foo.sol:Foo
				targetDir := filepath.Join(outDir, baseName)
				os.MkdirAll(targetDir, 0755)

				os.WriteFile(filepath.Join(targetDir, baseName+".abi"), []byte(contract.ABI), 0644)
				os.WriteFile(filepath.Join(targetDir, baseName+".bin"), []byte(contract.Bin), 0644)

				meta := DeploymentResult{
					Contract:   baseName,
					Source:     src,
					Network:    "devnet",
					Address:    "0xABCDEF123456...", // Placeholder
					TxHash:     "0xDEADBEEF...",     // Placeholder
					DeployedAt: time.Now().UTC(),
				}
				data, _ := json.MarshalIndent(meta, "", "  ")
				os.WriteFile(filepath.Join(targetDir, baseName+".deploy.json"), data, 0644)
				fmt.Println("💾 Written output to:", targetDir)
			}
		}
	},
}

func GetContractCommand() *cobra.Command {
	contractDeployCmd.Flags().String("src", "", "Path to the Solidity contract source file")
	contractDeployCmd.Flags().String("out", "", "Directory to output compiled contract artifacts")
	contractDeployCmd.MarkFlagRequired("src")
	contractCmd.AddCommand(contractDeployCmd)
	return contractCmd
}
