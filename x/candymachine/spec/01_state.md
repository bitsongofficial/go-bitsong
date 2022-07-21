# State

## Candy machine configuration

| Parameter             | Description                                               |
| :-------------------- | :-------------------------------------------------------- |
| collectionId          | Id of collection to mint                                  |
| price                 | price of nft to pay for minting an nft                    |
| number                | The number of items in the Candy Machine                  |
| treasury              | wallet to receive payments                                |
| denom                 | denom to be used for payment                              |
| goLiveDate            | Timestamp when minting is allowed                         |
| endSettings           | Describes mint end condition - by time or by minted count |
| whitelistMintSettings | Whitelist management rules for the machine                |
| metadataBaseUrl       | Base url of metadata for the collection                   |

When a collection owner creates candy machine, ownership of collection is sent to candy machine, and it can be returned only after candy machine end condition meets.

The tokens spent on minting nfts, is sent to the treasury account.

Note: How do we handle rare nfts on the collection with different price?
