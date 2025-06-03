package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Smart contract related operations",
}

var deployCmd = &cobra.Command{
	Use:   "deploy [path to .sol]",
	Short: "Compile and deploy a smart contract",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		solFile := args[0]
		fmt.Println("🚀 Compiling contract:", solFile)

		// compile via solc
		compile := exec.Command("solc", "--bin", "--abi", solFile, "-o", "build/")
		output, err := compile.CombinedOutput()
		if err != nil {
			fmt.Println("❌ Compile error:", err)
			fmt.Println(string(output))
			return
		}
		fmt.Println("✅ Compilation complete.")

		// TODO: Read bin+abi, send tx via Go-ethereum using your wallet logic
	},
}

func GetContractCommand() *cobra.Command {
	contractCmd.AddCommand(deployCmd)
	return contractCmd
}
