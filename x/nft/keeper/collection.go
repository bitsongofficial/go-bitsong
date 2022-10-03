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
	bz := store.Get(types.CollectionKey(id))
	if bz == nil {
		return types.Collection{}, sdkerrors.Wrapf(types.ErrCollectionDoesNotExist, "collection: %d does not exist", id)
	}
	collection := types.Collection{}
	k.cdc.MustUnmarshal(bz, &collection)
	return collection, nil
}

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) {
	bz := k.cdc.MustMarshal(&collection)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CollectionKey(collection.Id), bz)
}

func (k Keeper) GetAllCollections(ctx sdk.Context) []types.Collection {
	store := ctx.KVStore(k.storeKey)

	collections := []types.Collection{}
	it := sdk.KVStorePrefixIterator(store, types.PrefixCollection)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		collection := types.Collection{}
		k.cdc.MustUnmarshal(it.Value(), &collection)
		collections = append(collections, collection)
	}
	return collections
}
