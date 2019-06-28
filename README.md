<p align="center" background="black"><img src="bitsong-logo.png" width="398"></p>

[![License: APACHE](https://img.shields.io/badge/License-APACHE-yellow.svg)](https://github.com/BitSongOfficial/go-bitsong/blob/master/LICENSE)

This repository hosts `BitSong`, the first implementation of the BitSong Blockchain.

**BitSong** is a new music streaming platform based on [Tendermint](https://github.com/tendermint/tendermint) consensus algorithm and the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) toolkits. Please make sure you study these projects as well if you are not already familiar.

**BitSong** is a project dedicated to musicians and their fans, which aims to overcome the bureaucratic and economic obstacles within this industry and reward artists and users for simply using the platform.

On the **BitSong** platform you (artist) will be able to produce songs in which an advertiser can attach advertisements and users can access from any device. Funds via the Bitsong token `$BTSG` will be credited to the artist wallet immediately and they will be able to withdraw or convert as they see fit.

**Artists** need no longer to wait several months before a record label sends various reports, they can check the progress in real time directly within the Wallet.

_NOTE: This is alpha software. Please contact us if you aim to run it in production._

## BitSong Testnet

To run a full-node for the testnet of BitSong Blockchain, first [install `bitsongd`](./docs/installation.md), then follow [the guide](./docs/join-testnet.md).

For status updates and genesis file, see the [networks repo](https://github.com/BitSongOfficial/networks).

## Resources
- [Official Website](https://bitsong.io)

### Community
- [Telegram Channel (English)](https://t.me/BitSongOfficial)
- [Facebook](https://www.facebook.com/BitSongOfficial)
- [Twitter](https://twitter.com/BitSongOfficial)
- [Medium](https://medium.com/@BitSongOfficial)
- [Reddit](https://www.reddit.com/r/bitsong/)
- [BitcoinTalk ANN](https://bitcointalk.org/index.php?topic=2850943)
- [Linkedin](https://www.linkedin.com/company/bitsong)
- [Instagram](https://www.instagram.com/bitsong_official/)

## License

APACHE License

## Versioning

### SemVer

BitSong uses [SemVer](http://semver.org/) to determine when and how the version changes.
According to SemVer, anything in the public API can change at any time before version 1.0.0

To provide some stability to BitSong users in these 0.X.X days, the MINOR version is used
to signal breaking changes across a subset of the total public API. This subset includes all
interfaces exposed to other processes, but does not include the in-process Go APIs.
