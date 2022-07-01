# `merkledrop`

## Abstract

This document specifies the _merkledrop_ module of the BitSong chain.

The _merkledrop_ module enables the BitSong chain to support the ability to simply create custom token airdrops on chain. In this sense, any user can realize a custom airdrop for his community. 

Fantokens owners can use this module to airdrop their new minted fantokens to their fans. Similarly, also tokenholders are enabled to create airdrop for any occasion (e.g., as the result of a competition).

## Merkledrop

Based on the _Merkle Tree_ data structure, the merkledrops are very performant airdrops for the BitSong community. 

It is possible to identify each _merkledrop_ through its `id` and the module allows the users involved to autonomously claim tokens in the temporal window between the `StartHeight` and the `EndHeight` blocks.

Thanks to the _merkledrop_ module, users on BitSong can:

- manage _merkledrop_, create and claim them;
- build applications that use the _merkledrops_ API to create completely new custom airdrops.

Features that may be added in the future are described in Future Improvements.

## Table of Contents

1. **[Concepts](01_concepts.md)**
   - [Merkle Drop](01_concepts.md#Merkledrop)
   - [Merkle Tree](01_concepts.md#Merkle-tree)
   - [Verification process](01_concepts.md#Verification-process)
2. **[State](02_state.md)**
   - [Params](02_state.md#Params)
     <!--
     State Transitions
     -->
     <!--
     Keeper
     -->
3. **[Messages](03_messages.md)**
   - [MsgCreate](03_messages.md#MsgCreate)
   - [MsgClaim](03_messages.md#MsgClaim)
     <!--
     Begin-Block
     -->
4. **[End Block](04_end_block.md)**
   - [Withdraw](04_end_block.md#Withdraw)
   - [Delete completed merkledrop](04_end_block.md#Delete-completed-merkledrop)
5. **[Events](05_events.md)**
   - [EventCreate](04_events.md#EventCreate)
   - [EventClaim](04_events.md#EventClaim)
   - [EventWithdraw](04_events.md#EventWithdraw)
6. **[Parameters](05_parameters.md)**
   <!--
   Test Cases
   -->
   <!--
   Benchmarks
   -->
7. **[Client](07_client.md)**   
8. **[Future Improvements](08_future_improvements.md)**
