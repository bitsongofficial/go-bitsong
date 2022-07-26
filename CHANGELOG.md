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

## [v0.11.0] -2022-07-01

* (fantoken) introduce the [fantoken module](./x/fantoken/spec)
* (merkledrop) introduce the [merkledrop module](./x/merkledrop/spec)
* (app) bump [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) to [v0.45.6](https://github.com/cosmos/cosmos-sdk/tree/v0.45.6)
* (app) bump [ibc](https://github.com/cosmos/ibc-go) to [v3.0.0](https://github.com/cosmos/ibc-go/tree/v3.0.0)
* (app) bump [tendermint](https://github.com/tendermint/tendermint) to [v0.34.19](https://github.com/tendermint/tendermint/tree/v0.34.19)
* (app) bump [packet-forward-middleware](https://github.com/strangelove-ventures/packet-forward-middleware) to [v2.1.1](github.com/strangelove-ventures/packet-forward-middleware/tree/v2.1.1)
* (app) update swagger to reflect new modules
* (app) small fixs Makefile