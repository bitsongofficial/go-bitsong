# `cadence`

## Abstract

This document specifies the internal `x/cadence` module of Bitsong Network. This is a fork of the `x/clock` module from Juno.

The `x/cadence` module allows specific contracts to be executed at the end of every block. This allows the smart contract to perform actions that may need to happen every block or at set block intervals.

By using this module, your application can remove the headache of external whitelisted bots and instead depend on the chain itself for constant executions.

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Contract Integration](03_integration.md)**
