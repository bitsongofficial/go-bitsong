package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

func (k Keeper) GetNFTsByOwner(ctx sdk.Context, owner sdk.AccAddress) []types.NFT {
	store := ctx.KVStore(k.storeKey)

	nfts := []types.NFT{}
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixNFTByOwner, owner...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		id := string(it.Value())
		nft, err := k.GetNFTById(ctx, id)
		if err != nil {
			panic(err)
		}

		nfts = append(nfts, nft)
	}
	return nfts
}

func (k Keeper) GetCollectionNfts(ctx sdk.Context, collectionId uint64) []types.NFT {
	store := ctx.KVStore(k.storeKey)

	nfts := []types.NFT{}
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixNFT, sdk.Uint64ToBigEndian(collectionId)...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		nft := types.NFT{}
		k.cdc.MustUnmarshal(it.Value(), &nft)
		nfts = append(nfts, nft)
	}
	return nfts
}

func (k Keeper) GetNFTById(ctx sdk.Context, id string) (types.NFT, error) {
	store := ctx.KVStore(k.storeKey)

	if !types.IsValidNftId(id) {
		return types.NFT{}, types.ErrInvalidNftId
	}

	bz := store.Get(append(types.PrefixNFT, types.NftIdToBytes(id)...))
	if bz == nil {
		return types.NFT{}, sdkerrors.Wrapf(types.ErrNFTDoesNotExist, "nft: %s does not exist", id)
	}
	nft := types.NFT{}
	k.cdc.MustUnmarshal(bz, &nft)
	return nft, nil
}

func (k Keeper) SetNFT(ctx sdk.Context, nft types.NFT) {
	// check if previous NFT exists and delete
	if oldNft, err := k.GetNFTById(ctx, nft.Id()); err == nil {
		k.DeleteNFT(ctx, oldNft)
	}

	idBz := nft.IdBytes()
	bz := k.cdc.MustMarshal(&nft)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixNFT, idBz...), bz)

	owner, err := sdk.AccAddressFromBech32(nft.Owner)
	if err != nil {
		panic(err)
	}
	store.Set(append(append(types.PrefixNFTByOwner, owner...), idBz...), []byte(nft.Id()))
}

func (k Keeper) DeleteNFT(ctx sdk.Context, nft types.NFT) {
	idBz := nft.IdBytes()
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.PrefixNFT, idBz...))

	owner, err := sdk.AccAddressFromBech32(nft.Owner)
	if err != nil {
		panic(err)
	}
	store.Delete(append(append(types.PrefixNFTByOwner, owner...), idBz...))
}

func (k Keeper) GetAllNFTs(ctx sdk.Context) []types.NFT {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.PrefixNFT)
	defer it.Close()

	allNFTs := []types.NFT{}
	for ; it.Valid(); it.Next() {
		var nft types.NFT
		k.cdc.MustUnmarshal(it.Value(), &nft)

		allNFTs = append(allNFTs, nft)
	}

	return allNFTs
}

func (k Keeper) PayNftIssueFee(ctx sdk.Context, sender sdk.AccAddress) error {
	fee := k.GetParamSet(ctx).IssuePrice
	if fee.IsPositive() {
		feeCoins := sdk.Coins{fee}
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, feeCoins)
		if err != nil {
			return err
		}
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, feeCoins)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) CreateNFT(ctx sdk.Context, msg *types.MsgCreateNFT) (uint64, string, error) {

	collection, err := k.GetCollectionById(ctx, msg.CollId)
	if err != nil {
		return 0, "", err
	}

	if collection.UpdateAuthority != msg.Sender {
		return 0, "", types.ErrNotEnoughPermission
	}

	// create metadata
	metadataId := k.GetLastMetadataId(ctx) + 1
	k.SetLastMetadataId(ctx, metadataId)
	msg.Metadata.Id = metadataId
	for index := range msg.Metadata.Creators {
		msg.Metadata.Creators[index].Verified = false
	}
	if msg.Metadata.MasterEdition != nil {
		msg.Metadata.MasterEdition.Supply = 1
		if msg.Metadata.MasterEdition.MaxSupply < 1 {
			msg.Metadata.MasterEdition.MaxSupply = 1
		}
	}
	k.SetMetadata(ctx, msg.Metadata)
	ctx.EventManager().EmitTypedEvent(&types.EventMetadataCreation{
		Creator:    msg.Sender,
		MetadataId: msg.Metadata.Id,
	})

	// burn fees before minting an nft
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return 0, "", err
	}
	err = k.PayNftIssueFee(ctx, sender)
	if err != nil {
		return 0, "", err
	}

	// create nft
	nft := types.NFT{
		Owner:      msg.Sender,
		CollId:     msg.CollId,
		MetadataId: metadataId,
		Seq:        0,
	}
	k.SetNFT(ctx, nft)
	ctx.EventManager().EmitTypedEvent(&types.EventNFTCreation{
		Creator: msg.Sender,
		NftId:   nft.Id(),
	})
	return metadataId, nft.Id(), nil
}

func (k Keeper) PrintEdition(ctx sdk.Context, msg *types.MsgPrintEdition) (string, error) {
	metadata, err := k.GetMetadataById(ctx, msg.MetadataId)
	if err != nil {
		return "", err
	}

	masterEditionNft := types.NFT{
		CollId:     msg.CollId,
		MetadataId: msg.MetadataId,
		Seq:        0,
	}

	if _, err := k.GetNFTById(ctx, masterEditionNft.Id()); err != nil {
		return "", types.ErrMasterEditionNftDoesNotExists
	}

	if metadata.MintAuthority != msg.Sender {
		return "", types.ErrNotEnoughPermission
	}

	if metadata.MasterEdition == nil {
		return "", types.ErrNotMasterEditionNft
	}

	if metadata.MasterEdition.MaxSupply <= metadata.MasterEdition.Supply {
		return "", types.ErrAlreadyReachedEditionMaxSupply
	}

	edition := metadata.MasterEdition.Supply
	metadata.MasterEdition.Supply += 1
	k.SetMetadata(ctx, metadata)

	// burn fees before minting an nft
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return "", err
	}

	err = k.PayNftIssueFee(ctx, sender)
	if err != nil {
		return "", err
	}

	// create nft
	nft := types.NFT{
		Owner:      msg.Owner,
		CollId:     msg.CollId,
		MetadataId: msg.MetadataId,
		Seq:        edition,
	}
	k.SetNFT(ctx, nft)
	ctx.EventManager().EmitTypedEvent(&types.EventNFTCreation{
		Creator: msg.Sender,
		NftId:   nft.Id(),
	})

	return nft.Id(), nil
}

func (k Keeper) TransferNFT(ctx sdk.Context, msg *types.MsgTransferNFT) error {
	nft, err := k.GetNFTById(ctx, msg.Id)
	if err != nil {
		return err
	}

	if nft.Owner != msg.Sender {
		return types.ErrNotNFTOwner
	}

	nft.Owner = msg.NewOwner
	k.SetNFT(ctx, nft)
	ctx.EventManager().EmitTypedEvent(&types.EventNFTTransfer{
		NftId:    msg.Id,
		Sender:   msg.Sender,
		Receiver: msg.NewOwner,
	})

	return nil
}
