# Messages

## MsgCreateCandyMachine

```protobuf
message MsgCreateCandyMachine {
  string sender = 1;
  bitsong.candymachine.v1beta1.CandyMachine machine = 2 [ (gogoproto.nullable) = false ];
}
```

`MsgCreateCandyMachine` is the message to create candy machine by collection owner.

Steps:

1. Pay candymachine creation fee if fee exists
2. Ensure collection is owned by `msg.Sender`
3. Transfer ownership of collection to module account
4. Create candymachine object from provided params
5. Allocate mintable metadata ids after shuffle operation
6. Emit event for candymachine creation

## MsgUpdateCandyMachine

```protobuf
message MsgUpdateCandyMachine {
  string sender = 1;
  bitsong.candymachine.v1beta1.CandyMachine machine = 2 [ (gogoproto.nullable) = false ];
}
```

`MsgUpdateCandyMachine` is the message to update candy machine by candymachine authority.

Steps:

1. Ensure `msg.Sender` is candymachine authority
2. Update candymachine object from provided params
3. Emit event for candymachine update

## MsgCloseCandyMachine

```protobuf
message MsgCloseCandyMachine {
  string sender = 1;
  uint64 coll_id = 2;
}
```

`MsgCloseCandyMachine` is the message to close candy machine by candymachine authority. Collection ownership is sent back to the candymachine authority.

Steps:

1. Ensure `msg.Sender` is candymachine authority
2. Delete candymachine object from the store
3. Remove mintable metadata ids allocated for the candy machine
4. Update authority of collection to `msg.Sender`
5. Emit event for candymachine close

## MsgMintNFT

```protobuf
message MsgMintNFT {
  string sender = 1;
  uint64 collection_id = 2;
  string name = 3;
}
```

`MsgMintNFT` is the message to mint an nft through live candymachine.

Steps:

1. Ensure collection is put on candymachine
2. Ensure candymachine passed live date
3. Pay nft mint fee via candy machine
4. Get shuffled metadata id and create new metadata
5. Add new nft with new metadata
6. Increase the number of nfts minted by the machine
7. If minted count pass the threshold value `MaxMint`, close candymachine
8. Otherwise, store updated candymachine into the storage
9. Emit event for minting nft via candymachine
10. Collect nft id and return

## MsgMintNFTs

```protobuf
message MsgMintNFTs {
  string sender = 1;
  uint64 collection_id = 2;
  string number = 3;
}
```

`MsgMintNFTs` is the message to mint multiple nfts on a single message.

Steps:

1. Iterate `number` times
2. Execute single `MsgMintNft` message
3. Collect nft ids and return
