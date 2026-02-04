# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!--
Types of changes:
- **Security** - security patches and vulnerability fixes
- **State-Breaking** - changes requiring coordinated upgrade
- **Features** - new functionality
- **Bug Fixes** - bug fixes
- **Dependencies** - dependency updates
- **Documentation** - documentation changes
- **CI** - continuous integration changes

Format: * (**scope**) Description (#PR or commit)
-->

## [Unreleased]

### Security

### State-Breaking

### Features

### Bug Fixes

### Dependencies

### Documentation

### CI

* (**release**) Use native arm64 runners, disable Windows builds (#300)

---

## [0.23.1] - 2026-02-04

### Security

* (**cometbft**) Bump CometBFT to v0.38.21 ([CSA-2026-001](https://github.com/bitsongofficial/internal-docs/blob/main/docs/operations/security/advisories/CSA-2026-001.md)) (#297)

### CI

* (**release**) Modernize release workflow with multi-arch builds and checksums (#298)
* (**changelog**) Add changelog enforcement and conventional commit validation (#298)
* (**changelog**) Fix shell quoting in GITHUB_OUTPUT writes (#299)

---

## [0.20.0]

### Features

* (**upgrade**) Add v020 upgrade handler

### Bug Fixes

* (**staking**) Fix delegators inability to claim rewards after v018 patch

---

## [0.19.0]

### Features

* (**ci**) Add e2e tests for wasm, packet-forward-middleware & polytone (ibc-callbacks)
* (**docker**) Add Makefile for managing Docker-related tasks
* (**tests**) Add CosmWasm contract validation in testing framework
* (**app**) Enhance configuration capabilities with additional parameters
* (**cli**) Add commands for verifying slashed delegators and calculating reward discrepancies
* (**localnet**) Add Makefile for managing local testnet environment
* (**scripts**) Add automation script for downloading Polytone contracts
* (**docker**) Improve docker commands

### Bug Fixes

* (**gov**) Register legacy gov messages
* (**app**) Register ModuleBasics GrpcGatewayRoutes

### State-Breaking

* Remove versioning of app go module
* Remove randomGenesisAccounts param from auth module registration

---

## [0.18.0]

### Features

* (**tests**) Add interchaintest package support
* (**ci**) Add CI support to build & release docker image
* (**ci**) Add CI support to run interchain tests
* (**scripts**) Add `scripts/test_node.sh` for basic node setup testing

### Dependencies

* (**tendermint**) Replace tendermint with cometbft
* (**wasmd**) Bump to v0.45.0
* (**cosmos-sdk**) Bump to v0.47.8
* (**ibc-go**) Bump to v7.4.0
* (**pfm**) Bump packet-forward-middleware to v7.1.3

### State-Breaking

* Bump required minimum Go version to v1.22
* (**merkledrop**) Remove `x/merkledrop` module

---

## [0.16.0]

### Bug Fixes

* (**pfm**) Patch for Packet Forward Middleware bug

---

## [0.15.0] - 2024-03-06

### Dependencies

* (**cosmos-sdk**) Update to v0.45.16
* (**ibc-go**) Upgrade to v4.4.2
* (**tendermint**) Replace with CometBFT
* (**cosmwasm**) Upgrade to v0.33.0
* (**pfm**) Replace strangelove-ventures/packet-forward-middleware with cosmos/ibc-apps/middleware/packet-forward-middleware

---

## [0.14.0] - 2023-02-07

### Bug Fixes

* (**authz**) Add Binary Codec support to MinValCommissionDecorator
* (**authz**) Add MinValCommissionDecorator test

---

## [0.13.0] - 2023-01-23

### Features

* (**cli**) Add `init-from-state` command for easy initialization from exported state

### Dependencies

* (**cosmos-sdk**) Update to v0.45.11
* (**ibc-go**) Upgrade to v3.3.1
* (**tendermint**) Upgrade to v0.34.24
* (**cosmwasm**) Integrate v0.29.2

---

## [0.11.0] - 2022-07-01

### Features

* (**fantoken**) Introduce the [fantoken module](./x/fantoken/spec)
* (**merkledrop**) Introduce the [merkledrop module](./x/merkledrop/spec)
* (**docs**) Update swagger to reflect new modules

### Dependencies

* (**cosmos-sdk**) Bump to v0.45.6
* (**ibc-go**) Bump to v3.0.0
* (**tendermint**) Bump to v0.34.19
* (**pfm**) Bump packet-forward-middleware to v2.1.1

---

[Unreleased]: https://github.com/bitsongofficial/go-bitsong/compare/v0.23.1...HEAD
[0.23.1]: https://github.com/bitsongofficial/go-bitsong/compare/v0.20.0...v0.23.1
[0.20.0]: https://github.com/bitsongofficial/go-bitsong/compare/v0.19.0...v0.20.0
[0.19.0]: https://github.com/bitsongofficial/go-bitsong/compare/v0.18.0...v0.19.0
[0.18.0]: https://github.com/bitsongofficial/go-bitsong/compare/v0.16.0...v0.18.0
[0.16.0]: https://github.com/bitsongofficial/go-bitsong/compare/v0.15.0...v0.16.0
[0.15.0]: https://github.com/bitsongofficial/go-bitsong/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/bitsongofficial/go-bitsong/compare/v0.13.0...v0.14.0
[0.13.0]: https://github.com/bitsongofficial/go-bitsong/compare/v0.11.0...v0.13.0
[0.11.0]: https://github.com/bitsongofficial/go-bitsong/releases/tag/v0.11.0
