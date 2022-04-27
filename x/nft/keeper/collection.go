package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetLastCollectionId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastCollectionId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastCollectionId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastCollectionId, idBz)
}

func (k Keeper) GetCollectionById(ctx sdk.Context, id uint64) (types.Collection, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixMetadata, sdk.Uint64ToBigEndian(id)...))
	if bz == nil {
		return types.Collection{}, sdkerrors.Wrapf(types.ErrCollectionDoesNotExist, "metadata: %d does not exist", id)
	}
	collection := types.Collection{}
	k.cdc.MustUnmarshal(bz, &collection)
	return collection, nil
}

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) {
	idBz := sdk.Uint64ToBigEndian(collection.Id)
	bz := k.cdc.MustMarshal(&collection)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixCollection, idBz...), bz)
}

func (k Keeper) SetCollectionNftRecord(ctx sdk.Context, collectionId uint64, nftId uint64) {
	collectionIdBz := sdk.Uint64ToBigEndian(collectionId)
	nftIdBz := sdk.Uint64ToBigEndian(nftId)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(append(types.PrefixCollectionRecord, collectionIdBz...), nftIdBz...), sdk.Uint64ToBigEndian(nftId))
}

func (k Keeper) GetCollectionNftRecords(ctx sdk.Context, collectionId uint64) []uint64 {
	store := ctx.KVStore(k.storeKey)

	nftIds := []uint64{}
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixCollectionRecord, sdk.Uint64ToBigEndian(collectionId)...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		nftIds = append(nftIds, id)
	}
	return nftIds
}
