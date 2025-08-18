<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes.
"State Machine Breaking" for breaking the AppState

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog
## [v0.23.0]
## Added
## Fixed
### Removed 
### Improvements 
### Depreceated
### State Breaking
## [v0.22.0]
## Added
## Fixed
### Removed 
### Improvements 
### Depreceated
### State Breaking
## [v0.21.0]
## Added
## Fixed
### Removed 
### Improvements 
### Depreceated
### State Breaking
## [v0.20.0]

## Added
- v020 upgrade handler 

## Fixed
- certain delegators inability to claim rewards after v018 patch 
### Removed 
### Improvements 
### Depreceated
### State Breaking

## [v0.19.0]
### Added 
- new ci tests for wasm,packetforwardmiddleware & polytone (ibc-callbacks)
- Introduced a Makefile for managing Docker-related tasks, including commands for building various Docker images.
- Added functionality for validating CosmWasm contracts in a new testing framework.
- Enhanced configuration capabilities for the Bitsong application with additional parameters.
- Introduced commands for verifying slashed delegators and calculating discrepancies in delegator rewards.
- Added a new Makefile for managing a local testnet environment for the Bitsong blockchain.
- Introduced a new script for automating the downloading of Polyone contract.
- Improved docker commands

### Fixed
- Registered legacy gov msgs
- Register ModuleBasics GrpcGatewayRoutes

### Removed 
- versioning of app go module 
- removed randomGenesisAccounts as param on new apps auth module registration

## [v0.18.x]
### Added
- Interchaintest package support added
- New CI support to build & release docker image
- New CI support to run interchain tests
- New script at `scripts/test_node.sh` that is a basic script to test setting up and starting a node.
### Improvements 
- Improved Makefile cli top-level command scripting
- Replaced tendermint with cometbft
- Bumped wasmd to `v0.45.0`
- Bumped Cosmos-SDK to `v0.47.8`
- Bumped IBC-Go to `v7.4.0`
- Bumped Paket-Forward-Middleware to `v7.1.3`
- Reformatted app test suite

### Depreceated
- remove `x/merkledrop` module

### State Breaking 
- Bumped required minimum Go version to `v1.22`

## [v0.16.0]
### Bug fixes 
- patch for v0.16.0 that fixed Packet Forward Middleware bug.

## [v0.15.0] - 2024-03-06
- Updated Cosmos-sdk to v0.45.16 for improved stability and security
- Upgraded ibc-go to v4.4.2 for enhanced interoperability between different blockchain networks
- Replaced Tendermint with CometBFT
- Upgraded Cosmwasm to v0.33.0 for advanced smart contract functionality
- Replaced strangelove-ventures/packet-forward-middleware with cosmos/ibc-apps/middleware/packet-forward-middleware

## [v0.14.0] - 2023-02-07
- fix(authz): Add Binary Codec support to MinValCommissionDecorator
- fix(authz): Add MinValCommissionDecorator test

## [v0.13.0] - 2023-01-23
- Updated Cosmos-sdk to v0.45.11 for improved stability and security
- Upgraded ibc-go to v3.3.1 for enhanced interoperability between different blockchain networks
- Tendermint upgraded to v0.34.24 for better performance and bug fixes
- Integrated Cosmwasm v0.29.2 for advanced smart contract functionality
- Added new command init-from-state which allows for easy initialization of private validator, p2p, genesis, and application configuration, as well as replacement of exported state.

## [v0.11.0] -2022-07-01

* (fantoken) introduce the [fantoken module](./x/fantoken/spec)
* (merkledrop) introduce the [merkledrop module](./x/merkledrop/spec)
* (app) bump [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) to [v0.45.6](https://github.com/cosmos/cosmos-sdk/tree/v0.45.6)
* (app) bump [ibc](https://github.com/cosmos/ibc-go) to [v3.0.0](https://github.com/cosmos/ibc-go/tree/v3.0.0)
* (app) bump [tendermint](https://github.com/tendermint/tendermint) to [v0.34.19](https://github.com/tendermint/tendermint/tree/v0.34.19)
* (app) bump [packet-forward-middleware](https://github.com/strangelove-ventures/packet-forward-middleware) to [v2.1.1](github.com/strangelove-ventures/packet-forward-middleware/tree/v2.1.1)
* (app) update swagger to reflect new modules
* (app) small fixs Makefile