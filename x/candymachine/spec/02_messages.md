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
5. Emit event for candymachine creation

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
3. Update authority of collection to `msg.Sender`
4. Emit event for candymachine close

## MsgMintNFTResponse

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
3. Mint nft from module account with candymachine parameters and nft name passed as `msg.Name`
4. Transfer ownership of nft to `msg.Sender`
5. Increase the number of nfts minted by the machine
6. If end settings is by minted count and if minted count pass the threshold value on EndSettings, close candymachine
7. Otherwise, store updated candymachine into the storage
8. Emit event for minting nft via candymachine
