# What is BitSong?

**BitSong** is a new music streaming platform based on [Tendermint](https://github.com/tendermint/tendermint) consensus algorithm and the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) toolkits. Please make sure you study these projects as well if you are not already familiar.

**BitSong** is a project dedicated to musicians and their fans, which aims to overcome the bureaucratic and economic obstacles within this industry and reward artists and users for simply using the platform.

On the **BitSong** platform you (artist) will be able to produce songs in which an advertiser can attach advertisements and users can access from any device. Funds via the Bitsong token `$BTSG` will be credited to the artist wallet immediately and they will be able to withdraw or convert as they see fit.

**Artists** need no longer to wait several months before a record label sends various reports, they can check the progress in real time directly within the Wallet.

`go-bitsong` is the name of the Cosmos SDK application for the Cosmos Hub. It comes with 2 main entrypoints:

- `bitsongd`: The BitSong Daemon, runs a full-node of the `bitsong` application.
- `bitsongcli`: The BitSong command-line interface, which enables interaction with a BitSong full-node.

`bitsong` is built on the Cosmos SDK using the following modules:

- `x/auth`: Accounts and signatures.
- `x/bank`: Token transfers.
- `x/staking`: Staking logic.
- `x/mint`: Inflation logic.
- `x/distribution`: Fee distribution logic.
- `x/slashing`: Slashing logic.
- `x/gov`: Governance logic.
- `x/ibc`: Inter-blockchain transfers.
- `x/params`: Handles app-level parameters.
- `x/artist`: Artist logic.
- `x/song`: Song logic.

>About the Cosmos Hub: The Cosmos Hub is the first Hub to be launched in the Cosmos Network. The role of a Hub is to facilitate transfers between blockchains. If a blockchain connects to a Hub via IBC, it automatically gains access to all the other blockchains that are connected to it.


>BitSong is a public Proof-of-Stake chain. Its staking token is called the BTSG.