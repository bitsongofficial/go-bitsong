# Events

## Messages

### MsgCreateAuction

| Type                                           | Attribute Key | Attribute Value                       |
| :--------------------------------------------- | :------------ | :------------------------------------ |
| bitsong.marketplace.v1beta1.EventCreateAuction | creator       | {creator}                             |
| bitsong.marketplace.v1beta1.EventCreateAuction | auction_id    | {auction_id}                          |
| bitsong.nft.v1beta1.EventNFTTransfer           | nft_id        | {nft_id}                              |
| bitsong.nft.v1beta1.EventNFTTransfer           | sender        | {sender}                              |
| bitsong.nft.v1beta1.EventNFTTransfer           | receiver      | {receiver}                            |
| message                                        | action        | /bitsong.marketplace.MsgCreateAuction |

### MsgSetAuctionAuthority

| Type                                                 | Attribute Key | Attribute Value                             |
| :--------------------------------------------------- | :------------ | :------------------------------------------ |
| bitsong.marketplace.v1beta1.EventSetAuctionAuthority | auction_id    | {auction_id}                                |
| bitsong.marketplace.v1beta1.EventSetAuctionAuthority | authority     | {authority}                                 |
| message                                              | action        | /bitsong.marketplace.MsgSetAuctionAuthority |

### MsgStartAuction

| Type                                          | Attribute Key | Attribute Value                      |
| :-------------------------------------------- | :------------ | :----------------------------------- |
| bitsong.marketplace.v1beta1.EventStartAuction | auction_id    | {auction_id}                         |
| message                                       | action        | /bitsong.marketplace.MsgStartAuction |

### MsgEndAuction

| Type                                        | Attribute Key | Attribute Value                    |
| :------------------------------------------ | :------------ | :--------------------------------- |
| bitsong.marketplace.v1beta1.EventEndAuction | auction_id    | {auction_id}                       |
| message                                     | action        | /bitsong.marketplace.MsgEndAuction |

### MsgPlaceBid

| Type                                      | Attribute Key | Attribute Value                  |
| :---------------------------------------- | :------------ | :------------------------------- |
| bitsong.marketplace.v1beta1.EventPlaceBid | bidder        | {bidder}                         |
| bitsong.marketplace.v1beta1.EventPlaceBid | auction_id    | {auction_id}                     |
| coin_received                             | receiver      | {receiver}                       |
| coin_received                             | amount        | {amount}                         |
| coin_spent                                | spender       | {spender}                        |
| coin_spent                                | amount        | {amount}                         |
| message                                   | action        | /bitsong.marketplace.MsgPlaceBid |
| message                                   | sender        | {sender}                         |
| transfer                                  | recipient     | {recipient}                      |
| transfer                                  | sender        | {sender}                         |
| transfer                                  | amount        | {amount}                         |

### MsgCancelBid

| Type                                       | Attribute Key | Attribute Value                   |
| :----------------------------------------- | :------------ | :-------------------------------- |
| bitsong.marketplace.v1beta1.EventCancelBid | bidder        | {bidder}                          |
| bitsong.marketplace.v1beta1.EventCancelBid | auction_id    | {auction_id}                      |
| coin_received                              | receiver      | {receiver}                        |
| coin_received                              | amount        | {amount}                          |
| coin_spent                                 | spender       | {spender}                         |
| coin_spent                                 | amount        | {amount}                          |
| message                                    | action        | /bitsong.marketplace.MsgCancelBid |
| message                                    | sender        | {sender}                          |
| transfer                                   | recipient     | {recipient}                       |
| transfer                                   | sender        | {sender}                          |
| transfer                                   | amount        | {amount}                          |

### MsgClaimBid

| Type                                      | Attribute Key | Attribute Value                  |
| :---------------------------------------- | :------------ | :------------------------------- |
| bitsong.marketplace.v1beta1.EventClaimBid | bidder        | {bidder}                         |
| bitsong.marketplace.v1beta1.EventClaimBid | auction_id    | {auction_id}                     |
| coin_received                             | receiver      | {receiver}                       |
| coin_received                             | amount        | {amount}                         |
| coin_spent                                | spender       | {spender}                        |
| coin_spent                                | amount        | {amount}                         |
| message                                   | action        | /bitsong.marketplace.MsgClaimBid |
| message                                   | sender        | {sender}                         |
| transfer                                  | recipient     | {recipient}                      |
| transfer                                  | sender        | {sender}                         |
| transfer                                  | amount        | {amount}                         |

## Endblocker

| Type                                          | Attribute Key | Attribute Value |
| :-------------------------------------------- | :------------ | :-------------- |
| []bitsong.marketplace.v1beta1.EventEndAuction | auction_id    | {auction_id}    |
