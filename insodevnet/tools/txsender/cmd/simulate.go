package cmd

import (
	"fmt"
	"math/big"
	"time"

	"github.com/insoblok/flatgas/insodevnet/tools/txsender/cmd/simulate"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate repeated or random transactions",
	Long:  `Simulate transaction activity on the chain for testing or demonstration purposes.`,
}

var simulateFundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Simulate fund transfers",
	Run: func(cmd *cobra.Command, args []string) {
		base, _ := cmd.Flags().GetString("base")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		amountStr, _ := cmd.Flags().GetString("amount")
		rpcURL, _ := cmd.Flags().GetString("rpc")
		fromPassword, _ := cmd.Flags().GetString("password")
		count, _ := cmd.Flags().GetInt("count")
		interval, _ := cmd.Flags().GetDuration("interval")

		amount := new(big.Int)
		_, ok := amount.SetString(amountStr, 10)
		if !ok {
			fmt.Printf("❌ Invalid amount: %s\n", amountStr)
			return
		}

		fmt.Printf(
			"🚀 Starting simulation: sending %s wei from %s to %s times %d rpc %s\n",
			amount.String(),
			from,
			to,
			count,
			rpcURL,
		)
		for i := 1; i <= count; i++ {
			fmt.Printf("🔁 Tx %d/%d: sending funds...\n", i, count)

			err := simulate.SendFunds(base, from, fromPassword, to, amount, rpcURL)
			if err != nil {
				fmt.Printf("❌ Tx %d failed: %v\n", i, err)
			} else {
				fmt.Printf("✅ Tx %d sent\n", i)
			}

			if i < count {
				time.Sleep(interval)
			}
		}
		fmt.Println("🎉 Simulation done.")
	},
}
var simulateScenarioSubCmd = &cobra.Command{
	Use:   "scenario",
	Short: "Run a transaction simulation scenario from a JSON specification",
	Run: func(cmd *cobra.Command, args []string) {
		specPath, _ := cmd.Flags().GetString("spec")
		err := simulate.DoScenario(specPath)
		if err != nil {
			fmt.Println("❌ failed to run scenario: %v", err)
		}
		fmt.Println("🎉 Simulation done.")
	},
}

func init() {
	initFundSubCmd()
	initScenarioSubCmd()
}

func initFundSubCmd() {
	simulateCmd.AddCommand(simulateFundCmd)
	simulateFundCmd.Flags().String("base", ".", "Base path to flatgas root")
	simulateFundCmd.Flags().String("from", "", "Sender alias")
	simulateFundCmd.Flags().String("password", "", "Password for sender account")
	simulateFundCmd.Flags().String("to", "", "Recipient alias or address")
	simulateFundCmd.Flags().String("amount", "0", "Amount in wei")
	simulateFundCmd.Flags().Int("count", 1, "Number of transactions to send")
	simulateFundCmd.Flags().String("rpc", "http://localhost:8545", "RPC endpoint URL or alias")
	simulateFundCmd.Flags().Duration("interval", 2*time.Second, "Interval between transactions")

	_ = simulateFundCmd.MarkFlagRequired("base")
	_ = simulateFundCmd.MarkFlagRequired("from")
	_ = simulateFundCmd.MarkFlagRequired("to")
	_ = simulateFundCmd.MarkFlagRequired("amount")
	_ = simulateFundCmd.MarkFlagRequired("rpc")
	_ = simulateFundCmd.MarkFlagRequired("password")
}
func initScenarioSubCmd() {
	simulateCmd.AddCommand(simulateScenarioSubCmd)
	simulateScenarioSubCmd.Flags().String("spec", "", "Path to the JSON scenario file defining simulation parameters and transaction behavior")
	_ = simulateScenarioSubCmd.MarkFlagRequired("spec")
}

func GetSimulateCommand() *cobra.Command {
	return simulateCmd
}
