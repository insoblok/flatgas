package main

import (
	"fmt"
	"time"
)

// Transaction represents a simple Flatgas transaction
type Transaction struct {
	ID          string
	Timestamp   time.Time
	GasUsed     uint64
	IsEmergency bool
}

// Mempool is a simple FIFO queue of transactions
type Mempool struct {
	Queue []Transaction
}

// Add a transaction to the mempool
func (m *Mempool) AddTx(tx Transaction) {
	m.Queue = append(m.Queue, tx)
}

// Pop transactions up to a gas limit
func (m *Mempool) PopTxs(maxGas uint64) ([]Transaction, uint64) {
	var included []Transaction
	var gasUsed uint64

	for len(m.Queue) > 0 && gasUsed+m.Queue[0].GasUsed <= maxGas {
		tx := m.Queue[0]
		m.Queue = m.Queue[1:]
		included = append(included, tx)
		gasUsed += tx.GasUsed
	}
	return included, gasUsed
}

// Block represents a mined block
type Block struct {
	Transactions []Transaction
	TotalGasUsed uint64
	TotalFees    uint64
}

const FixedGasPrice uint64 = 1000 // arbitrary unit, for showoff

func main() {
	mempool := Mempool{}

	// Simulate submitting transactions
	for i := 0; i < 10; i++ {
		tx := Transaction{
			ID:          fmt.Sprintf("tx%d", i),
			Timestamp:   time.Now(),
			GasUsed:     21000,
			IsEmergency: false,
		}
		mempool.AddTx(tx)
		time.Sleep(100 * time.Millisecond) // simulate small delay
	}

	// Simulate block production
	fmt.Println("Producing a block...")
	includedTxs, gasUsed := mempool.PopTxs(1000000) // 1M gas limit
	totalFees := gasUsed * FixedGasPrice

	block := Block{
		Transactions: includedTxs,
		TotalGasUsed: gasUsed,
		TotalFees:    totalFees,
	}

	fmt.Printf("Block created with %d transactions, total gas used: %d, total fees collected: %d units\n", len(block.Transactions), block.TotalGasUsed, block.TotalFees)

	// List included transactions
	for _, tx := range block.Transactions {
		fmt.Printf("Included %s (Gas: %d)\n", tx.ID, tx.GasUsed)
	}
}
