# Fan Token Specification

## Abstract
This specification describes how to issue a new fantoken on the chain. A fantoken is a simple Fungible Token interface.
There are two keys to identify a fantoken. It is denom and symbol.
Both represents the fantoken but they are different.
Denom is unique and calculated by tendermint crypto hash with creator, symbol and name.
Symbol is defined by the user.

## Contents

1. **[State](01_state.md)**
    - [FanToken](01_state.md#FanToken)
    - [Params](01_state.md#Params)
2. **[Messages](02_messages.md)**
    - [MsgIssueFanToken](02_messages.md#MsgIssueFanToken)
    - [MsgEditFanToken](02_messages.md#MsgEditFanToken)
    - [MsgMintFanToken](02_messages.md#MsgMintFanToken)
    - [MsgBurnFanToken](02_messages.md#MsgBurnFanToken)
    - [MsgTransferFanTokenOwner](02_messages.md#MsgTransferFanTokenOwner)
3. **[Events](03_events.md)**
    - [MsgIssueFanToken](03_events.md#MsgIssueFanToken)
    - [MsgEditFanToken](03_events.md#MsgEditFanToken)
    - [MsgMintFanToken](03_events.md#MsgMintFanToken)
    - [MsgBurnFanToken](03_events.md#MsgBurnFanToken)
    - [MsgTransferFanTokenOwner](03_events.md#MsgTransferFanTokenOwner)
4. **[Parameters](04_params.md)**