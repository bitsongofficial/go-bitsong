# `fantoken`

## Abstract

This document specifies the _fantoken_ module of the BitSong chain.

The _fantoken_ module enables the BitSong chain to support fan tokens, allowing actors in the content creation industry to create their economy. In this sense, they can generate new ways to monetize their music and brand and provide a unique and innovative channel to engage with fans. Thanks to this module, players from the content creation universe can start minting their _fan tokens_ (which are **fungible tokens**) and listing them within a few minutes for low fees.

### An example: Fan tokens in the music Industry

In the music industry, for example, _fan tokens_ enable to empower a lot of different scenarios. For instance, it is possible to use them to crowdfund a tour or an album, or even to access exclusive content. The potential of such a system is very massive and, with these few examples, you can imagine what a contribution this tool can make to a world teeming with content creators.

### Fan tokens in BitSong

Based on the concept of the **ERC-20 Standard**, BitSong _fan tokens_ enable the user to a new way of **value exchanging**. Here, through tokens issued by a particular entity, the fans can deeply interact with their influencers or idols.

We can identify each _fan token_ through its `denom`.
Moreover, even if its `denom` allow the global identification of the token, each _fan token_ is also equipped with a `name` and a `symbol`, which helps in its recognition. The `name` and the `symbol` of a _fan token_, together with a `uri` and an `authority` (i.e., the address of the wallet which is able to manage those data) are part of the `metadata` of the _fan token_.

More specifically:

- **denom** is calculated by the tendermint crypto hash function through the *block height* of the transaction, the first *minter*, the *symbol*, and the *name*. For this reason, it is _unique_;
- **symbol** is defined by the user and can be any string matching the pattern `^[a-z0-9]{1,64}$`, so any lowercase string containing letters and digits with a length between 1 and 64 characters. _It cannot be empty_;
- **name**, on the other hand, is also defined by the user but it can be any string containing max 128 characters. _It can also be empty_.

Finally, thanks to the _fantoken_ module, users on BitSong can:

- manage _fan tokens_, issuing, minting, burning, and transferring them;
- build applications that use the _fan tokens_ API to create completely new and custom artists' economies.

Features that may be added in the future are described in Future Improvements.

## Table of Contents

1. **[Concepts](01_concepts.md)**
   - [Conventions](01_concepts.md#Conventions)
   - [Fan Token](01_concepts.md#Fan-token)
   - [Lifecycle of a fan token](01_concepts.md#Lifecycle-of-a-fan-token)
   - [Uniqueness of the denom](01_concepts.md#Uniqueness-of-the-denom)
2. **[State](02_state.md)**
   - [Params](02_state.md#Params)
   - [Fan Token](02_state.md#Token)
     <!--
     State Transitions
     -->
     <!--
     Keeper
     -->
3. **[Messages](03_messages.md)**
   - [MsgIssue](03_messages.md#MsgIssue)
   - [MsgDisableMint](03_messages.md#MsgDisableMint)
   - [MsgMint](03_messages.md#MsgMint)
   - [MsgBurn](03_messages.md#MsgBurn)
   - [MsgSetAuthority](03_messages.md#MsgSetAuthority)
   - [MsgSetMinter](03_messages.md#MsgSetMinter)
   - [MsgSetUri](03_messages.md#MsgSetUri)
     <!--
     Begin-Block
     -->
     <!--
     End-Block
     -->
4. **[Events](04_events.md)**
   - [EventIssue](04_events.md#EventIssue)
   - [EventDisableMint](04_events.md#EventDisableMint)
   - [EventMint](04_events.md#EventMint)
   - [EventBurn](04_events.md#EventBurn)
   - [EventSetAuthority](04_events.md#EventSetAuthority)
   - [EventSetMinter](04_events.md#EventSetMinter)
   - [EventSetUri](04_events.md#EventSetUri)
5. **[Parameters](05_parameters.md)**
   <!--
   Test Cases
   -->
   <!--
   Benchmarks
   -->
6. **[Future Improvements](06_future_improvements.md)**
