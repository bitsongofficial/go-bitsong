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
- cosmwasm: download or build binaries to reduce overall workspace size & remove prebuilt binaries

## Local Development

To test a ci workflow prior to pushing to a git repo, we make use of the [act](https://nektosact.com/):

```sh
# runs a specific workflow, need to specify the correct container-architecture if on arm64
act -W '.github/workflows/tests.yml' --container-architecture linux/amd64
```