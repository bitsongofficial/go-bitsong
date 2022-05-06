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
