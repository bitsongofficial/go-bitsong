## Bash Tests

This repo contains bash scripts for testing local Bitsong networks.

Here's the markdown table for those modules:

| Module | Description | Documentation |
|--------|-------------|---------------|
| [AA (Abstract Accounts)](./interchain-accounts/) | | [README](./interchain-accounts/) |
| [IBC-Hooks](./ibchook/) | | [README](./ibchook/) |
| [x/nft](./nft/) | | [README](./nft/) |
| [PFM - Packet Forward Middleware](./pfm/README) | | [README](./pfm/README) |
| [Polytone](./polytone/README) | | [README](./polytone/README) |
| [Staking-Hooks](./staking-hooks/README) | | [README](./staking-hooks/README) |
| [Upgrades](./upgrade/) | | [README](./upgrade/) |
<!-- | [Cadence](./cadence/) | | [README](./cadence/) | -->
 
## TODO

- refactor to remove redundancy of starting up nodes
- upgrades: improve downloading snapshot
- cosmwasm: download or build binaries to reduce overall workspace size & remove prebuilt binaries
- wire into cicd
