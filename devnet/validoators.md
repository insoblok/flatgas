# Validators in Flatgas

## 🧱 What Is a Validator?

A **validator** is a special type of account that:

* Has a **private key** and an **on-chain balance**.
* Runs node software capable of:

    * **Proposing** new blocks.
    * **Validating** and **attesting** to blocks created by others.
* Signs blocks using its private key, establishing authorship.

---

## 🔄 How Is a Validator Selected?

The validator selection method depends on the consensus mechanism. For Flatgas (early stage PoS), options include:

1. **Round-robin** — Sequential rotation.
2. **Weighted random** — Based on staked `inso` amount.
3. **Committee-based** — Epoch-driven validator groups.

For early `insotestnet`, a static list in `genesis.json` with round-robin logic is simplest.

---

## 🔍 Validator vs. Smart Contract

* A **smart contract** is passive: it only acts when **invoked** by a transaction.
* A **validator** includes and **executes** transactions that interact with contracts.

### Smart Contract Flow:

1. User submits transaction calling a smart contract.
2. Validator picks it up from mempool.
3. Executes contract logic (e.g., transfer, mint).
4. Records the result in a block.
5. Signs and broadcasts the block.

---

## ❌ Can a Block Be Rejected?

Yes. Any **full node** re-validates each block. If:

* State transitions are invalid
* Contract logic is tampered
* Signatures are wrong

Then the node **rejects the block**. Result:

* The network forks.
* Faulty validators may be **slashed** or removed (future feature).

---

## ✅ Summary

| Topic                        | Details                                           |
| ---------------------------- | ------------------------------------------------- |
| **Validator account**        | Has private key, balance; used for signing blocks |
| **Validator node**           | Runs consensus code, proposes/validates blocks    |
| **Block signing**            | Proves authorship, enables consensus              |
| **Selection method**         | Round-robin, stake-based, committee, etc.         |
| **Smart contract execution** | Performed by validators on tx inclusion           |
| **Block validation**         | All nodes re-execute & can reject invalid blocks  |

---

## 🔄 Upgradeable Validator Set

Best practice is to eventually allow:

* Adding/removing validators **dynamically**
* Managing this via **on-chain governance** or **protocol rules**

Until then, `genesis.json` defines the static validator list.

---

## 🚀 Transition to Production

In `insotestnet`, validators can be defined and managed manually.
Later, production `insonet` should:

* Have proper validator registration
* Incentivize staking
* Penalize misbehavior (slashing)
* Use upgradeable validator management
