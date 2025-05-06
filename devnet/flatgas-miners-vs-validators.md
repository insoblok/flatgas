
# Flatgas: Miners vs Validators

In Flatgas (as in post-Merge Ethereum), the traditional concept of "miners" has been replaced by **validators** due to the shift from Proof-of-Work (PoW) to Proof-of-Stake (PoS).

---

## 🛠️ Miners (Legacy - PoW)

- Used in Ethereum pre-Merge
- Solve computational puzzles (Proof-of-Work)
- Validate and include transactions in blocks
- Consume real-world energy (electricity)
- Receive mining rewards + transaction fees

> The `--mine` flag in Geth was used to enable mining in PoW mode.
> In Flatgas, this flag is now deprecated and has no effect.

---

## ✅ Validators (Current - PoS)

- Used in Flatgas and post-Merge Ethereum
- Propose and validate blocks without computational puzzles
- Must be authorized or staked to participate
- Receive transaction fees or protocol rewards
- Included in `genesis.json` or added later via governance

---

## ⚙️ Why Validators Are Important

- Maintain consensus and block finality
- Ensure network security
- Participate in governance (if protocol allows)
- Can be slashed for malicious behavior

---

## 📦 Genesis Configuration

Validators are often defined in the `genesis.json`:

```json
"extraData": "0x000000000000000000000000<validator1>000000000000000000000000<validator2>..."
```

- This bootstraps trusted validator identities for the network start.
- Useful for testnets or small PoA-style test environments.

---

## 🔄 Adding Validators Later

- Yes, possible via:
  - Governance mechanism (DAO, multisig, voting)
  - Smart contract validator registry
  - Manual reconfiguration + restart (less ideal)
- Best practice: use upgradeable validator set via protocol rules.

---

## 📝 Summary

| Role       | Old (PoW)             | New (PoS)               |
|------------|-----------------------|--------------------------|
| Name       | Miner                 | Validator                |
| Mechanism  | Solve puzzle (PoW)    | Stake & validate         |
| Power cost | High                  | Low                      |
| Rewards    | Mining + fees         | Fees or protocol rewards |
| Included   | N/A                   | Genesis or added later   |

---

Flatgas is **PoS-native**: no mining, just validators.
Stay focused on setting up robust validator sets and their upgrade paths.

