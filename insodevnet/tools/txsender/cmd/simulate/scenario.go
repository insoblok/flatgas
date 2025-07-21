package simulate

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"
)

type Scenario struct {
	RPC            string        `json:"rpc"`
	Base           string        `json:"base"`
	TxCount        int           `json:"txCount"`
	MinIntervalSec int           `json:"minIntervalSec"`
	MaxIntervalSec int           `json:"maxIntervalSec"`
	Actions        []ActionEntry `json:"actions"`
}

type ActionEntry struct {
	Type         string        `json:"type"` // currently only "fund", future: "contractCall", etc.
	FromAccounts []FromAccount `json:"fromAccounts"`
	ToAccounts   []string      `json:"toAccounts"`
	AmountWei    []string      `json:"amountWei"` // string to avoid issues with big.Int parsing
}

type FromAccount struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

func DoScenario(specPath string) error {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("❌ failed to read spec file: %w", err)
	}

	var scenario Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return fmt.Errorf("❌ failed to parse scenario: %w", err)
	}

	fmt.Println(scenario)
	fmt.Println()
	if len(scenario.Actions) == 0 {
		log.Fatalf("❌ no actions defined in scenario")
	}
	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= scenario.TxCount; i++ {
		fmt.Printf("🔁 Tx %d/%d: sending funds...\n", i, scenario.TxCount)

		action := scenario.Actions[rand.Intn(len(scenario.Actions))]

		if len(action.FromAccounts) == 0 {
			log.Fatalf("❌ no fromAccounts defined in action")
		}

		from := action.FromAccounts[rand.Intn(len(action.FromAccounts))]
		fmt.Printf("🔍 Chosen sender: %s (with password)\n", from.Alias)

		if len(action.ToAccounts) == 0 {
			log.Fatalf("❌ no toAccounts defined in action")
		}

		to := action.ToAccounts[rand.Intn(len(action.ToAccounts))]
		fmt.Printf("🎯 Chosen recipient: %s\n", to)

		if len(action.AmountWei) == 0 {
			log.Fatalf("❌ no amountWei options defined in action")
		}

		amountStr := action.AmountWei[rand.Intn(len(action.AmountWei))]
		amount := new(big.Int)
		_, ok := amount.SetString(amountStr, 10)
		if !ok {
			log.Fatalf("❌ failed to parse amount: %s", amountStr)
		}

		fmt.Printf("💰 Chosen amount (wei): %s\n", amount.String())

		fromResolved := MustResolve(scenario.Base, from.Alias)

		toResolved := MustResolve(scenario.Base, to)

		err := SendFunds(scenario.Base, fromResolved.Hex(), from.Password, toResolved.Hex(), amount, scenario.RPC)
		if err != nil {
			fmt.Printf("❌ Tx %d failed: %v\n", i, err)
		} else {
			fmt.Printf("✅ Tx %d sent\n", i)
		}

		if i < scenario.TxCount {
			minSec := scenario.MinIntervalSec
			maxSec := scenario.MaxIntervalSec
			if maxSec < minSec {
				log.Fatalf("❌ maxIntervalSec (%d) is less than minIntervalSec (%d)", maxSec, minSec)
			}

			randomSec := rand.Intn(maxSec-minSec+1) + minSec
			interval := time.Duration(randomSec) * time.Second
			fmt.Printf("⏱️ Waiting for %s before next tx\n", interval)
			time.Sleep(interval)
		}
	}

	return nil
}
