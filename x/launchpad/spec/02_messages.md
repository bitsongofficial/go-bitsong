# Messages

## MsgCreateLaunchPad

```protobuf
message MsgCreateLaunchPad {
  string sender = 1;
  bitsong.launchpad.v1beta1.LaunchPad pad = 2 [ (gogoproto.nullable) = false ];
}
```

`MsgCreateLaunchPad` is the message to create launchpad by collection owner.

Steps:

1. Pay launchpad creation fee if fee exists
2. Ensure collection is owned by `msg.Sender`
3. Transfer ownership of collection to module account
4. Create launchpad object from provided params
5. Allocate mintable metadata ids after shuffle operation
6. Emit event for launchpad creation

## MsgUpdateLaunchPad

```protobuf
message MsgUpdateLaunchPad {
  string sender = 1;
  bitsong.launchpad.v1beta1.LaunchPad pad = 2 [ (gogoproto.nullable) = false ];
}
```

`MsgUpdateLaunchPad` is the message to update launchpad by launchpad authority.

Steps:

1. Ensure `msg.Sender` is launchpad authority
2. Update launchpad object from provided params
3. Emit event for launchpad update

## MsgCloseLaunchPad

```protobuf
message MsgCloseLaunchPad {
  string sender = 1;
  uint64 coll_id = 2;
}
```

`MsgCloseLaunchPad` is the message to close launchpad by launchpad authority. Collection ownership is sent back to the launchpad authority.

Steps:

1. Ensure `msg.Sender` is launchpad authority
2. Delete launchpad object from the store
3. Remove mintable metadata ids allocated for the launchpad
4. Update authority of collection to `msg.Sender`
5. Emit event for launchpad close

## MsgMintNFT

```protobuf
message MsgMintNFT {
  string sender = 1;
  uint64 collection_id = 2;
  string name = 3;
}
```

`MsgMintNFT` is the message to mint an nft through live launchpad.

Steps:

1. Ensure collection is put on launchpad
2. Ensure launchpad passed live date
3. Pay nft mint fee via launchpad
4. Get shuffled metadata id and create new metadata
5. Add new nft with new metadata
6. Increase the number of nfts minted by the pad
7. If minted count pass the threshold value `MaxMint`, close launchpad
8. Otherwise, store updated launchpad into the storage
9. Emit event for minting nft via launchpad
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
