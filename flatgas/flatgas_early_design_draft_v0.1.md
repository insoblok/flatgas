# Flatgas Protocol: Early Design Draft (v0.1)

Prepared by: Iyad & Flatgas Core Team  
Date: April 2025

---

# Abstract

Flatgas is a next-generation Layer 1 blockchain designed with deterministic, stable transaction costs and a focus on fairness, neutrality, and protocol simplicity. This document outlines the early design principles, gas pricing model, governance framework, and transaction inclusion philosophy.

---

# Executive Summary

**Flatgas** is a next-generation Layer 1 blockchain protocol built on the principle of **Determinism over Speculation**.

The goal is to create a blockchain where transaction costs are simple, predictable, and immune to speculative manipulation. 
Flatgas is designed to:
- Use a **fixed gas unit cost** rather than a dynamic market-driven base fee.
- Provide a stable, user-friendly experience.
- Prioritize fairness, transparency, and censorship resistance.

**Native Token:** `inso`

---

# 1. Gas Pricing Model

| # | Aspect | Ethereum Now (EIP-1559) | Flatgas Vision | Analysis & Considerations |
|:--|:-------|:-------------------------|:---------------|:--------------------------|
| 1.1 | Gas Pricing Model | Dynamic baseFee, varies per block based on block gas usage. | Fixed unit gas price, predictable across all blocks. | Simpler UX; Risk if network congestion changes over time; might need future governance to adjust. |
| 2.1 | Who Sets Fees? | Protocol sets base fee algorithmically. | Protocol fixes gas price for set periods; updates via governed review. | Governance needed; Risk of governance capture; See Cost Governance section. |
| 3.1 | Miner/Validator Rewards | Miners get tips only (base fee burned). | Validators get full flat fee (no burn by default). | Strong validator incentive; no burn reduces deflationary pressure. |
| 3.2 | Base Fee Usage | Base fees are burned, tips go to miners. | Flat fee can be split: 100% to validator, or partial burn. | Option to keep some burn for scarcity; needs economic modeling. |
| 4.1 | Mempool Behavior | Priority sorted: higher tip -> higher chance to be mined. | First-Come-First-Serve (possible). | Fairer transaction processing; needs DOS protection mechanisms. |
| 4.2 | Transaction Replacement | Higher tip transactions can replace lower ones (Replace-By-Fee - RBF). | No RBF based on fees; possible explicit cancel/resend mechanisms. | Safer UX; needs clear stuck-tx handling. |
| 5.1 | Risk of Stuck Transactions | High if base fee rises after signing. | Very low - flat fees mean transactions either fit or don't. | Great UX; simplifies wallet logic. |
| 5.2 | Transaction Confirmation Speed | Variable based on fee bidding. | Stable confirmation speed, based only on block capacity. | Predictable, fair block fill behavior. |
| 6.1 | Fee Burn Mechanism | Base fee burned every transaction. | No automatic burn, or optional configurable burn percentage. | Full incentive to validators or partial scarcity for inso tokenomics. |
| 6.2 | Token (`inso`) Economics | ETH deflationary due to burns. | `inso` flexible: inflationary, deflationary, or stable. | Needs early economic policy decision. |

---

# 2. Flatgas Cost Governance Plan

## Summary
Flatgas fixes the gas unit price for predetermined periods (e.g., 1 year). The fee is only reviewed based on strict, deterministic conditions.

## Key Principles
- **Review Timing:** Only after minimum fixed periods (e.g., 12 months).
- **Trigger Conditions:** Must meet protocol-visible, measurable metrics.
- **Decision Mechanism:** Governance DAO or strict protocol-controlled logic.
- **Advance Notice:** 30–90 days before fee adjustments are activated.
- **No Validator Control:** Validators cannot set or influence gas prices independently.

## Benefits
- User trust and predictability.
- Reduced manipulation risk.
- Aligns with core Flatgas philosophy.

---

# 3. Transaction Inclusion Policy

## Default Rule: FIFO (First-In-First-Out)

- Transactions are included in blocks strictly based on order of arrival.
- No transaction reordering based on fee size (fees are fixed anyway).
- Protects users from miner greed and transaction censorship.

## Special Cases: Emergency Transactions

- Open topic: Should special "emergency" tx types exist (e.g., validator resignation)?
- Needs extremely strict, verifiable rules to avoid abuse.
- Possible options: protocol-defined emergency types, strict whitelist.

## Benefits
- Radical fairness.
- Strengthened network neutrality.
- Minimized miner manipulation.

---

# 4. Open Questions for Future Refinement

- How exactly to structure governance for fee review?
- How to define (or if to allow) emergency transaction lanes?
- Should `inso` token economics favor slight inflation, deflation, or hard supply cap?
- How to implement strong DOS protection without fee prioritization?

---

# End of Draft v0.1

---

*This document represents an early internal draft of Flatgas design principles, intended for discussion and refinement.*
