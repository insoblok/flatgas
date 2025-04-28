# Flatgas Transaction Simulator

## Overview

This is a simple Go-based simulator built to model the basic transaction handling behavior of the Flatgas Layer 1 protocol.

Flatgas is a blockchain designed with the principle of **Determinism over Speculation**:
- Fixed gas unit cost (no dynamic base fee adjustment).
- FIFO (First-In-First-Out) transaction queue (no fee-based transaction jumps).
- Transparent validator reward calculation.
- Simplified block production logic.

This simulator demonstrates:
- Submitting transactions to a FIFO mempool.
- Producing blocks based on a fixed gas limit.
- Calculating total block fees collected by validators.

---

## How to Run

1. Make sure you have Go installed (`go version`).
2. Clone this repository or copy the `main.go` file.
3. Run the simulator:

```bash
go run main.go
