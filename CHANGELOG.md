# Changelog

All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

## [unreleased]

---

## [0.24.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.23.0..v0.24.0) - 2025-08-29

### ⚙️ Miscellaneous

- remove icq support - ([6e30322](https://github.com/bitsongofficial/go-bitsong/commit/6e30322327e85dd6acbc9796429640bc60f91872)) - hard-nett
- [**breaking**]bump wasmd@v0.61.3, wasmvm@v3.0.2, cometbft@0.38.18, pfm@v10.1, wasm-light-client@v10, ibc-go@v10, ibc-hooks@v10 - ([3ad1785](https://github.com/bitsongofficial/go-bitsong/commit/3ad17852b7027bf69455e28b5441cc0ee0bbba9b)) - hard-nett
- remove the rest of stale upgradehandler imports - ([f0ed673](https://github.com/bitsongofficial/go-bitsong/commit/f0ed673c6c9b4142f430955d966829e51316b972)) - hard-nett
- re-wire default export cmd - ([a1a260c](https://github.com/bitsongofficial/go-bitsong/commit/a1a260c96a066f0c13a100e9b306bf317072b2ff)) - hard-nett
- re-add ibc+wasm ict tests (via polytone) - ([a39e677](https://github.com/bitsongofficial/go-bitsong/commit/a39e67729e8a049eb54433af3075999b44ae704f)) - hard-nett

### 🐛 Bug fixes
- remove x/cadence - ([62c6b84](https://github.com/bitsongofficial/go-bitsong/commit/62c6b8434cc20c97d6d287c72b385ebdaf0dd982)) - hard-nett
- add x/fantoken MsgServiceHandler - ([160f17a](https://github.com/bitsongofficial/go-bitsong/commit/160f17a4014769936848f9aebbd3f7685ab4601e)) - hard-nett
- wire in v024 upgradehandler - ([ad31d93](https://github.com/bitsongofficial/go-bitsong/commit/ad31d93cb56adac2caf9b84627d5097e98ac243c)) - hard-nett
- wire 08-wasm & 07-tendermint light clients into cilentKeeper, use pointers for our AppKeepers struct - ([d32f36e](https://github.com/bitsongofficial/go-bitsong/commit/d32f36e9068e10ec752aa6f15a1a2cc88649c4c3)) - hard-nett
- revert bankKeeper pointer for now, improve pointer use throughout app - ([0756c5c](https://github.com/bitsongofficial/go-bitsong/commit/0756c5c45a3ce20abce57be2884d815d97606df8)) - hard-nett
- revert to proper ics4Middleware keeper initialization sequence - ([54feebf](https://github.com/bitsongofficial/go-bitsong/commit/54feebfbf8432417714309dbae472ec92eb9aa74)) - hard-nett

---
## [0.23.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.22.0..v0.23.0) - 2025-06-23

### ⚙️ Miscellaneous

- v023 unit tests - ([915094e](https://github.com/bitsongofficial/go-bitsong/commit/915094e03aed4402a9d81ea95418d8e06885e031)) - hard-nett
- - resolve v018-v021 issues (#278) - ([cf7a58c](https://github.com/bitsongofficial/go-bitsong/commit/cf7a58c6f3c274e9153f0f907ebfa8b88e58cf97)) - Hard-Nett
- go mod tidy - ([028b79e](https://github.com/bitsongofficial/go-bitsong/commit/028b79e1e8fee997b67cdc63157e0d2db8f61c17)) - hard-nett
- fantoken basic tests - ([48fbae7](https://github.com/bitsongofficial/go-bitsong/commit/48fbae762d63209ddaaf8bb23ff25da004984963)) - hard-nett
- fix test addr - ([03379ea](https://github.com/bitsongofficial/go-bitsong/commit/03379ea25165c98a391838f178ca6350d5802528)) - hard-nett
- go mod tidy - ([6bcd4ec](https://github.com/bitsongofficial/go-bitsong/commit/6bcd4ec9f49c9cca49dce7de5fe83f0c34c3ba03)) - hard-nett
- set resolveDenom to false, remove v022 custom patch from init-from-mainnet - ([a899c13](https://github.com/bitsongofficial/go-bitsong/commit/a899c1388cd02f473e050accc07051f8880a7563)) - hard-nett

### 🐛 Bug fixes

- add cosmos.msg.v1.service to protos - ([24820a3](https://github.com/bitsongofficial/go-bitsong/commit/24820a3ebdad599503095ebf45ff1479bc5bebe7)) - hard-nett
- remove external community pool support - ([0645ca9](https://github.com/bitsongofficial/go-bitsong/commit/0645ca9307a172b864fecbf66d51b1f824aefc64)) - hard-nett
- set feepool to reflect updated balance - ([831d0ee](https://github.com/bitsongofficial/go-bitsong/commit/831d0ee582d651e1dbc27715c740e38ad01225bc)) - hard-nett

---
## [0.22.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.21.6..v0.22.0) - 2025-05-18

### ⚙️ Miscellaneous

- - resolve v018-v021 issues (#278) - ([cf7a58c](https://github.com/bitsongofficial/go-bitsong/commit/cf7a58c6f3c274e9153f0f907ebfa8b88e58cf97)) - Hard-Nett

---
## [0.21.6](https://github.com/bitsongofficial/go-bitsong/compare/v0.21.5..v0.21.6) - 2025-03-13

### ⚙️ Miscellaneous

- upgrade ibc to 8.7.0 (#281) - ([fbdc845](https://github.com/bitsongofficial/go-bitsong/commit/fbdc84594b0a65d2bbe1da4381f71932efa7efed)) - Angelo RC

---
## [0.20.2](https://github.com/bitsongofficial/go-bitsong/compare/v0.20.1..v0.20.2) - 2024-12-27

### ⚙️ Miscellaneous

- V0.20.2 (#252) - ([e49371a](https://github.com/bitsongofficial/go-bitsong/commit/e49371a6876f650fc908ee376337606b2f57f3b5)) - Hard-Nett

---
## [0.18.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.17.0..v0.18.0) - 2024-11-27

### ⚙️ Miscellaneous

- cosmos-sdk@v0.47, wasmd@v0.46, ibc-go@v7, (#235) - ([50b4082](https://github.com/bitsongofficial/go-bitsong/commit/50b4082736a68cdde098cf36edd7c7a70d9fdae6)) - Hard-Nett

---
## [0.17.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.16.0..v0.17.0) - 2024-08-21

### 🐛 Bug fixes

- Remove IBCFeeKeeper and replace with ChannelKeeper in app.go (#234) - ([6caaf5f](https://github.com/bitsongofficial/go-bitsong/commit/6caaf5fdba8e7ce41e8a9d44654c141f85c9c38f)) - Angelo RC

---
## [0.16.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.15.0..v0.16.0) - 2024-08-09

### ⚙️ Miscellaneous

- pfm v4 (#232) - ([cf24787](https://github.com/bitsongofficial/go-bitsong/commit/cf247879efb67393a5e230244146456ae5638238)) - Angelo RC
- add node service handler (#233) - ([608415a](https://github.com/bitsongofficial/go-bitsong/commit/608415adfd0d5fedb081c369268461541b3fe4e4)) - Sergey

---
## [0.14.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.14.0-alpha.1..v0.14.0) - 2023-02-07

### ⚙️ Miscellaneous

- add v014 upgrade handler - ([e5847ac](https://github.com/bitsongofficial/go-bitsong/commit/e5847ac0b4f0319811d62140f4e438d33163ad0d)) - angelorc

### 🐛 Bug fixes

- **(authz)** Add MinValCommissionDecorator test - ([62fe771](https://github.com/bitsongofficial/go-bitsong/commit/62fe771e0ef98c2849d75495d9fb6b2df52d2489)) - angelorc
- **(authz)** Add MinValCommissionDecorator test - ([a6b3d76](https://github.com/bitsongofficial/go-bitsong/commit/a6b3d76bd00650dae0fc6c9323723ce4ffa547a5)) - angelorc

---
## [0.14.0-alpha.1](https://github.com/bitsongofficial/go-bitsong/compare/v0.13.0..v0.14.0-alpha.1) - 2023-02-04

### 🐛 Bug fixes

- **(authz)** Add Binary Codec support to MinValCommissionDecorator - ([8b3cc2e](https://github.com/bitsongofficial/go-bitsong/commit/8b3cc2eeb0ef48624ebf4517c6f6dce285798f5a)) - angelorc

---
## [0.13.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.12.1.1..v0.13.0) - 2023-01-24

### 🐛 Bug fixes

- go.sum - ([ecf29a7](https://github.com/bitsongofficial/go-bitsong/commit/ecf29a742bebfd9c037a9aeadb4e4ac85d4bcd4b)) - angelorc

---
## [0.12.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.11.0..v0.12.0) - 2022-10-14

### ⚙️ Miscellaneous

- fix incorrect license (#148) - ([604e886](https://github.com/bitsongofficial/go-bitsong/commit/604e886d8cc31df68d5c5abdd39a0836c2a3c420)) - Christophe Camel
- bump cosmos-sdk and go - ([66608bf](https://github.com/bitsongofficial/go-bitsong/commit/66608bfaa0ed8bed4c08827b3d2dee5cc437bf3c)) - angelorc

### 🐛 Bug fixes

- **(docker)** fix and push the new docker images - ([a7d1a18](https://github.com/bitsongofficial/go-bitsong/commit/a7d1a1893b8050c1679889cd507f59a56c197b51)) - angelorc

---
## [0.11.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.11.0-rc0..v0.11.0) - 2022-07-01

### ⚙️ Miscellaneous

- **(app)** default fee for merkledrop and fantoken is 1_000_000_000ubtsg - ([f92cb0e](https://github.com/bitsongofficial/go-bitsong/commit/f92cb0eced47db4bc9fdba21a16a2e1a1b9ebf85)) - angelorc

---
## [0.11.0-rc0](https://github.com/bitsongofficial/go-bitsong/compare/v0.10.0..v0.11.0-rc0) - 2022-07-01

### ⚙️ Miscellaneous

- **(app)** change golang version to gorelease - ([820b9d9](https://github.com/bitsongofficial/go-bitsong/commit/820b9d95c0cb87d2f39ade62a667d45ee5d3cc38)) - angelorc

---
## [0.10.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.8.0..v0.10.0) - 2022-01-31

### ⚙️ Miscellaneous

- change module minter for prop 6 - ([f84631f](https://github.com/bitsongofficial/go-bitsong/commit/f84631fd0721084772c1b4c77063af2c318bf3e0)) - angelorc
- change module minter for prop 6 - ([290f06c](https://github.com/bitsongofficial/go-bitsong/commit/290f06c9562b640dd4419f4170ee876b6149773b)) - angelorc
- ante handle val min commission - ([e4ffcf3](https://github.com/bitsongofficial/go-bitsong/commit/e4ffcf35e3e72b3257430a33a9401fc512db77a4)) - angelorc

---
## [0.8.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.8.0-rc1..v0.8.0) - 2021-10-05

### ⚙️ Miscellaneous

- add consensus_params.evidence.max_bytes in migrate genesis cmd - ([df4b8e4](https://github.com/bitsongofficial/go-bitsong/commit/df4b8e4f9ad8b43371c24d8e1a7746e9e1c20f12)) - Angelo
- swagger - ([995a19f](https://github.com/bitsongofficial/go-bitsong/commit/995a19f33d39d7c4156569215e2f6b9e057589ec)) - Angelo
- workflow - ([860f893](https://github.com/bitsongofficial/go-bitsong/commit/860f893138fffdd5c8812299a45c2b4d58a8bf14)) - Angelo

---
## [0.8.0-dev](https://github.com/bitsongofficial/go-bitsong/compare/v0.7.0-rc1..v0.8.0-dev) - 2021-06-23

### ⚙️ Miscellaneous

- add fantoken module - ([617af5e](https://github.com/bitsongofficial/go-bitsong/commit/617af5e374c27187f612d8c505d9566bca103b05)) - sarawut
- app name, tx cli, package path - ([a1f43ff](https://github.com/bitsongofficial/go-bitsong/commit/a1f43ffdfbcc1093d793a5ef1dccd2528765d2df)) - Angelo
- genesis overwrite - ([4a343ac](https://github.com/bitsongofficial/go-bitsong/commit/4a343ac0aa5543987a164f5c18b5dbe424181be1)) - sarawut
- fix mint fantoken amount #56 - ([555b7db](https://github.com/bitsongofficial/go-bitsong/commit/555b7dbb73ecc10347b8241d02b08901867fe158)) - sarawut

---
## [0.7.0-rc1](https://github.com/bitsongofficial/go-bitsong/compare/v0.6.0-beta1..v0.7.0-rc1) - 2021-03-20

### ⚙️ Miscellaneous

- default crisis GenesisState - ([276f739](https://github.com/bitsongofficial/go-bitsong/commit/276f739d2193d031a32e418723b1f136fd2cd112)) - Angelo

---
## [0.5.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.4.0..v0.5.0) - 2020-08-06

### ⚙️ Miscellaneous

- query shares - ([ca226c0](https://github.com/bitsongofficial/go-bitsong/commit/ca226c0ca156754432a98d9368c06e3ec5fa43a3)) - Angelo

### 🚀 New features

- add/remove share to track - ([00c29e3](https://github.com/bitsongofficial/go-bitsong/commit/00c29e33c36a2a8b531366d8cf3452ff43fefe55)) - Angelo

---
## [0.4.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.3.1..v0.4.0) - 2020-08-04

### ⚙️ Miscellaneous

- mint & burn to content module - ([00f43cc](https://github.com/bitsongofficial/go-bitsong/commit/00f43cc848ee8d840438e3a856c080889e1e3c95)) - Angelo
- content module (add/stream/download) - ([d998273](https://github.com/bitsongofficial/go-bitsong/commit/d998273587ec7156790f749d136692f8bb5fadc4)) - Angelo
- content module (add/stream/download) - ([4f376f7](https://github.com/bitsongofficial/go-bitsong/commit/4f376f786ce312b2b01cd40b912599a64bef7576)) - Angelo
- content module add rights holders - ([517f544](https://github.com/bitsongofficial/go-bitsong/commit/517f5442cd128032e1f739a982553e67649dd38a)) - Angelo
- integration player module - ([59f9e24](https://github.com/bitsongofficial/go-bitsong/commit/59f9e24d26e4f4321fa5720bb2f5cbe65f689fe0)) - Angelo
- tokenize rights - ([3112161](https://github.com/bitsongofficial/go-bitsong/commit/31121617f7ba50d701d0298af84006519b05de43)) - Angelo

---
## [0.3.0-dev-1](https://github.com/bitsongofficial/go-bitsong/compare/v0.2.1..v0.3.0-dev-1) - 2019-12-05

### ⚙️ Miscellaneous

- Add song pool parameter - ([103b5f3](https://github.com/bitsongofficial/go-bitsong/commit/103b5f34761fe505fe289034a890dbfdd0d0ff94)) - Angelo
- Renamed SongsPool to Rewards - ([987bff4](https://github.com/bitsongofficial/go-bitsong/commit/987bff4fea52183e1167d58bb43f53b91af24cbf)) - Angelo
- Added Play Struct - ([158ed63](https://github.com/bitsongofficial/go-bitsong/commit/158ed635f2bac1e1ad7c35a7f308900bc2b0bc39)) - Angelo
- Add get and set Play to keeper - ([d8f39fb](https://github.com/bitsongofficial/go-bitsong/commit/d8f39fba6b172da6a792955e629fa3bf189620fa)) - Angelo
- Refactor, added initial tests - ([b6730d2](https://github.com/bitsongofficial/go-bitsong/commit/b6730d28a35ad3c9485aa46fbffa35a7e87330ec)) - Angelo
- Continuos implementations and tests - ([706c96b](https://github.com/bitsongofficial/go-bitsong/commit/706c96bc8f0bcb0377324cbc7744f48f09e7eb94)) - Angelo
- Add GetUserPower to Keeper - ([4fe9092](https://github.com/bitsongofficial/go-bitsong/commit/4fe90928b31c51228bef441e66602b41d8050f79)) - Angelo
- Breaking changes! - ([65cc4bc](https://github.com/bitsongofficial/go-bitsong/commit/65cc4bc729302a90b8ae19671d8f0ae8b2c2d150)) - Angelo
- fix genesis params - ([886df73](https://github.com/bitsongofficial/go-bitsong/commit/886df733217f1b398824bedb2d26cf63b3f4920a)) - Angelo
- Added play get/set to keeper - ([ac3c8f4](https://github.com/bitsongofficial/go-bitsong/commit/ac3c8f47c75f6571e5fde15e6d795e47d376c0a0)) - Angelo
- added SavePlay to keeper - ([9523519](https://github.com/bitsongofficial/go-bitsong/commit/952351972a2107a78a5ff4f91c9e43a89d517894)) - Angelo
- added tests - ([0b1a4e9](https://github.com/bitsongofficial/go-bitsong/commit/0b1a4e9ad3abe8f1a431ae457ef49416acb5549d)) - Angelo
- GetAllPlays keeper - ([44debf0](https://github.com/bitsongofficial/go-bitsong/commit/44debf0f6a77eb732824f379af29e606d5e60b55)) - Angelo
- Fix cli tx PublishTrack - ([a4a5f41](https://github.com/bitsongofficial/go-bitsong/commit/a4a5f419a89a7011c58fe19453521bae70dbfb49)) - Angelo
- Removed duplicated events - ([cd0e96b](https://github.com/bitsongofficial/go-bitsong/commit/cd0e96b6c7fad0f8fedacede90c076f4668b404c)) - Angelo
- Fix cli tx PlayTrack - ([4b4cb49](https://github.com/bitsongofficial/go-bitsong/commit/4b4cb49e95c8f8d74199fdde7dd1d416d15ecf70)) - Angelo
- Renamed SavePlay to PlayTrack - ([feb85f6](https://github.com/bitsongofficial/go-bitsong/commit/feb85f64eda47f662976abe87125ee55a95fcce9)) - Angelo
- Extended Cosmos-SDK Distribution Module - ([61719c0](https://github.com/bitsongofficial/go-bitsong/commit/61719c0a4818641ee1c0d5a7cff6aa58e9c3e0ef)) - Angelo
- Play Pool Rewards!!! - ([e2a9c77](https://github.com/bitsongofficial/go-bitsong/commit/e2a9c77745b097711db5c95a0e96f9a5d9723b61)) - Angelo
- Initial abci implementation - ([0a6c8fa](https://github.com/bitsongofficial/go-bitsong/commit/0a6c8fac9a4aa40e63897db3802446d01643fe72)) - Angelo
- ABCI initial logic to pay rewards - ([c3760e6](https://github.com/bitsongofficial/go-bitsong/commit/c3760e61ce329ce68ce9988903a5986652d8c28e)) - Angelo
- Tests improvements - ([d92aef7](https://github.com/bitsongofficial/go-bitsong/commit/d92aef75b289a3df022070d1fe39ac7951819ae7)) - Angelo
- test abci reward - ([62a2480](https://github.com/bitsongofficial/go-bitsong/commit/62a2480ca7d48f6539288a12967ef5aa9b80a64e)) - Angelo
- Continuous abci implementation - ([f636c01](https://github.com/bitsongofficial/go-bitsong/commit/f636c015f14bbf20bbfbbb14d0659f4caa87cfa2)) - Angelo
- AllocateTokens completed (need tests) - ([900a457](https://github.com/bitsongofficial/go-bitsong/commit/900a457b5cebd8bf0ce559c05f83b0500f48342e)) - Angelo
- abci (pay rewards) wip - ([220913f](https://github.com/bitsongofficial/go-bitsong/commit/220913ff3e967dfe47658085779cd7b115d48f4f)) - Angelo
- Continous implementations - ([170b6fd](https://github.com/bitsongofficial/go-bitsong/commit/170b6fda6be9c23b21078fc3f803f76e56e6dbdf)) - Angelo
- fix pay fee - ([d79ef57](https://github.com/bitsongofficial/go-bitsong/commit/d79ef57488cf9d5614c2259289905abebd99c670)) - Angelo

---
## [0.1.0](https://github.com/bitsongofficial/go-bitsong/compare/v0.0.2..v0.1.0) - 2019-07-10

### ⚙️ Miscellaneous

- bitsong-testnet-1 - ([bd6d3ec](https://github.com/bitsongofficial/go-bitsong/commit/bd6d3eca2a4ec453e07fd1ff5370d4e34022dd2b)) - Angelo

### 🚀 New features

- gentx for alessiotreglia - ([2080ee6](https://github.com/bitsongofficial/go-bitsong/commit/2080ee6389c4ef68364807f6942170acd04ec2e3)) - Alessio Treglia
- gentx for validator-center - ([e8a0e61](https://github.com/bitsongofficial/go-bitsong/commit/e8a0e613bdc0ca4b20ef20651b555bb3371a267e)) - zartus2019
- gentx for ManaComm - ([03b733c](https://github.com/bitsongofficial/go-bitsong/commit/03b733cb9c24d96ed916d597c5474ca5f737747b)) - ManaComm
- gentx for BitGlad - ([9269406](https://github.com/bitsongofficial/go-bitsong/commit/92694069aec55cb3365746ee4a8c1cc55e1e4126)) - BitGlad1
- gentx for BitAngel - ([3ead0e0](https://github.com/bitsongofficial/go-bitsong/commit/3ead0e0d7e3120bd55e90b74b940f83fa5730086)) - root
- gentx for <your-moniker> - ([c44dde6](https://github.com/bitsongofficial/go-bitsong/commit/c44dde66793a084fab1c94c04c157051daa8ac4a)) - Ubuntu
- gentx for mikmik - ([9d36be1](https://github.com/bitsongofficial/go-bitsong/commit/9d36be1884b0df771e8ec7985d4909d2ada86045)) - Alessio Treglia
- gentx for ondin - ([ae604de](https://github.com/bitsongofficial/go-bitsong/commit/ae604de755d4c68485d5154798d51ce3cf592629)) - Ondin777
- gentx for darkeyesz - ([160ce90](https://github.com/bitsongofficial/go-bitsong/commit/160ce901cefd7ea78a83129f6b749b711228a9aa)) - root
- gentx for anamix - ([a6a1d99](https://github.com/bitsongofficial/go-bitsong/commit/a6a1d992af0ca89032e48e2858198f4f89a42163)) - dbpatty
- gentx for redpenguin - ([9b0df8f](https://github.com/bitsongofficial/go-bitsong/commit/9b0df8f8a594c91a5f40ba443dba72ceaa58e18f)) - redpenguin-validator
- gentx for UbikCapital - ([a96a1d8](https://github.com/bitsongofficial/go-bitsong/commit/a96a1d8a3bb2a6350c6a13eae076a69da69660b8)) - root
- gentx for Simply VC - ([29625b2](https://github.com/bitsongofficial/go-bitsong/commit/29625b2b7151213c35ea4329877b3d59b12e2272)) - Simply VC
- gentx for genesislab - ([dab8b0b](https://github.com/bitsongofficial/go-bitsong/commit/dab8b0bd15c08c2076ceb07af842368fb1315dcd)) - i7495
- gentx for lxgn - ([7b5546c](https://github.com/bitsongofficial/go-bitsong/commit/7b5546ccc361792be7c4b8abdfced975ec240292)) - root
- gentx for bitcat - ([6f8b053](https://github.com/bitsongofficial/go-bitsong/commit/6f8b0530d7844a354fc6866c900a5bbb397d2473)) - Bit Cat
- gentx for GioStake - ([c4530fb](https://github.com/bitsongofficial/go-bitsong/commit/c4530fb0cf0af1b383431dbf258c3e5925020a27)) - ilgio
- gentx for gota - ([4a9d3a6](https://github.com/bitsongofficial/go-bitsong/commit/4a9d3a6c94a56484077093f1bde9389c7105f662)) - Ubuntu
- gentx for funky-validator - ([5256222](https://github.com/bitsongofficial/go-bitsong/commit/52562229544477d9e7f469d98b6bec1cce0743f1)) - zenfunkpanda
- gentx for zx - ([06d5157](https://github.com/bitsongofficial/go-bitsong/commit/06d5157958214c9e7c188e85e81b8f3555e43778)) - Juan Leni
- gentx for sebytza05 - ([9b45c45](https://github.com/bitsongofficial/go-bitsong/commit/9b45c4540f932111fea12464b168b44b0af62dd6)) - dirmansebastian
- gentx for commercio.network - ([a4b4e33](https://github.com/bitsongofficial/go-bitsong/commit/a4b4e33348bdc59b03db55c1e693e89cfc8c0736)) - marcotradenet
- gentx for mintonium - ([90d6362](https://github.com/bitsongofficial/go-bitsong/commit/90d6362b3637b881a095df8ab23989d5aab6c649)) - jaygaga
- gentx for forbole - ([4f27aac](https://github.com/bitsongofficial/go-bitsong/commit/4f27aac20e81176bf7cc63ca54e4ec2a74c07a1e)) - Kwun Yeung
- gentx for POS-Bakerz - ([3e7b06b](https://github.com/bitsongofficial/go-bitsong/commit/3e7b06b64f18c2ccc9cea06f1f129906a28c8292)) - root

---
## [0.0.1] - 2019-07-02

### ⚙️ Miscellaneous

- readme, license - ([e6e4c24](https://github.com/bitsongofficial/go-bitsong/commit/e6e4c24800244da9494cab86f2634d827719174e)) - Angelo
- Song module - ([69efcb1](https://github.com/bitsongofficial/go-bitsong/commit/69efcb1515927e1ea6b85012cb72c1f25c30af5f)) - Angelo
- Song module - ([b35ab6e](https://github.com/bitsongofficial/go-bitsong/commit/b35ab6e3c69ff7c072134a04ad068704b3eee836)) - Angelo
- Song module - ([080a1b5](https://github.com/bitsongofficial/go-bitsong/commit/080a1b5bbc7ec157ffcccdd422e3b73255ab910a)) - Angelo
- Song module - ([f5831e7](https://github.com/bitsongofficial/go-bitsong/commit/f5831e7d4b7c8ffb78537c9128749e7c021f6b71)) - Angelo
- added Content, TotalReward. RedistribuitionSplitRate - ([0da2c84](https://github.com/bitsongofficial/go-bitsong/commit/0da2c8435784eefab81429ee20701580c583c46a)) - Angelo
- minor fix - ([7aeb263](https://github.com/bitsongofficial/go-bitsong/commit/7aeb2637e0615c34a12fe7168482811d7b77f62a)) - Angelo
- WIP - ([53edb13](https://github.com/bitsongofficial/go-bitsong/commit/53edb131442eaa9914f07797d7d285b09a93f9ca)) - Angelo
- initial implementation - ([5343eb0](https://github.com/bitsongofficial/go-bitsong/commit/5343eb0a8386f91e66b92789b7f8f529fecc8ab5)) - Angelo
- command line - ([ef96c3d](https://github.com/bitsongofficial/go-bitsong/commit/ef96c3df036128bbb82b6dd19fb2622d23f081ff)) - Angelo

### 🚀 New features

- init command - ([a4872fa](https://github.com/bitsongofficial/go-bitsong/commit/a4872fa9671c1d31178cec393a96329f7c30e439)) - Angelo

<!-- generated by git-cliff -->
