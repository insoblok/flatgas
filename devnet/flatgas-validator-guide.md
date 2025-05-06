
# Flatgas Validator Primer

This document provides a concise summary of what validators are, their purpose in the Flatgas network, how they are configured in the genesis file, and how to manage them over time.

---

## 🧾 What is a Validator?

In Flatgas (like in Ethereum and other blockchains), a **validator** is a node that participates in block production and consensus. Validators are essential to:

- **Propose and validate new blocks**
- **Secure the network**
- **Earn rewards for honest behavior**
- **Face penalties or slashing for misbehavior (in Proof-of-Stake setups)**

---

## 🎯 Purpose of Validators

Validators serve as the backbone of consensus. Depending on the consensus mechanism:

- In **Clique (PoA)**: A fixed list of validators take turns to propose blocks.
- In **Proof-of-Stake**: Validators are selected based on stake and can be dynamically rotated.

---

## 📜 Including Validators in `genesis.json`

When using Clique (PoA), validators must be included in the `extraData` field of `genesis.json`:

```json
"extraData": "0x[padding] + [validator1_address] + [validator2_address] + ... + [padding]"
```

Each validator must have:
- An externally owned account (EOA) with a known address
- The same genesis file as other nodes

---

## ➕ Adding Validators Later

- In **Clique**, validator changes are done **on-chain** using **votes** from existing validators.
- In **PoS**, validator set is dynamic based on stake and protocol rules.
- **They do NOT need to be in the initial genesis if the protocol supports dynamic changes.**

---

## 🧪 Validator Considerations in Testnets

In a testnet, you are free to:
- Try static validator sets first (simpler setup)
- Experiment with validator change mechanics
- Validate block times, reward schedules, consensus behavior
- Practice simulating validator failures and recoveries

---

## 🏁 Recommendations for Flatgas Testnet

- Start with 1–3 validators included in genesis
- Use static peers or `static-nodes.json` for deterministic connectivity
- Use the `--nodekey` flag for fixed identity
- Track validator addresses and enodes in versioned config files

---

## 🧱 Genesis Location (Suggested)

```
flatgas/networks/insotestnet/genesis.json
```

--- 

This file is a living guide. As your testnet evolves, update it to reflect what you learn and standardize your validator setup practices.

