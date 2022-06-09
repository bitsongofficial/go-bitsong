package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

func (k Keeper) GetLastNftId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastNftId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastNftId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastNftId, idBz)
}

func (k Keeper) GetNFTsByOwner(ctx sdk.Context, owner sdk.AccAddress) []types.NFT {
	store := ctx.KVStore(k.storeKey)

	nfts := []types.NFT{}
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixNFTByOwner, owner...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		nft, err := k.GetNFTById(ctx, id)
		if err != nil {
			panic(err)
		}

		nfts = append(nfts, nft)
	}
	return nfts
}

func (k Keeper) GetNFTById(ctx sdk.Context, id uint64) (types.NFT, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixNFT, sdk.Uint64ToBigEndian(id)...))
	if bz == nil {
		return types.NFT{}, sdkerrors.Wrapf(types.ErrNFTDoesNotExist, "nft: %d does not exist", id)
	}
	nft := types.NFT{}
	k.cdc.MustUnmarshal(bz, &nft)
	return nft, nil
}

func (k Keeper) SetNFT(ctx sdk.Context, nft types.NFT) {
	// check if previous NFT exists and delete
	if oldNft, err := k.GetNFTById(ctx, nft.Id); err == nil {
		k.DeleteNFT(ctx, oldNft)
	}

	idBz := sdk.Uint64ToBigEndian(nft.Id)
	bz := k.cdc.MustMarshal(&nft)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixNFT, idBz...), bz)

	owner, err := sdk.AccAddressFromBech32(nft.Owner)
	if err != nil {
		panic(err)
	}
	store.Set(append(append(types.PrefixNFTByOwner, owner...), idBz...), idBz)
}

func (k Keeper) DeleteNFT(ctx sdk.Context, nft types.NFT) {
	idBz := sdk.Uint64ToBigEndian(nft.Id)
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

func (k Keeper) PrintEdition(ctx sdk.Context, msg *types.MsgPrintEdition) (uint64, error) {
	metadata, err := k.GetMetadataById(ctx, msg.MetadataId)
	if err != nil {
		return 0, err
	}

	if metadata.UpdateAuthority != msg.Sender {
		return 0, types.ErrNotEnoughPermission
	}

	if metadata.MasterEdition == nil {
		return 0, types.ErrNotMasterEditionNft
	}

	if metadata.MasterEdition.MaxSupply <= metadata.MasterEdition.Supply {
		return 0, types.ErrAlreadyReachedEditionMaxSupply
	}

	edition := metadata.MasterEdition.Supply
	metadata.MasterEdition.Supply += 1
	k.SetMetadata(ctx, metadata)

	// burn fees before minting an nft
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return 0, err
	}

	err = k.PayNftIssueFee(ctx, sender)
	if err != nil {
		return 0, err
	}

	// create nft
	nftId := k.GetLastNftId(ctx) + 1
	k.SetLastNftId(ctx, nftId)
	nft := types.NFT{
		Id:         nftId,
		Owner:      msg.Sender,
		MetadataId: msg.MetadataId,
		Edition:    edition,
	}
	k.SetNFT(ctx, nft)
	ctx.EventManager().EmitTypedEvent(&types.EventNFTCreation{
		Creator: msg.Sender,
		NftId:   nftId,
	})

	return nftId, nil
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
