package keeper

import (
	"math/rand"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetMintableMetadataIds(ctx sdk.Context, collId uint64, metadataIds []uint64) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, append(types.PrefixMintableMetadataIds, sdk.Uint64ToBigEndian(collId)...))

	for index, metadataId := range metadataIds {
		prefixStore.Set(sdk.Uint64ToBigEndian(uint64(index)), sdk.Uint64ToBigEndian(metadataId))
	}
}

func (k Keeper) TakeOutRandomMintableMetadataId(ctx sdk.Context, collId uint64, maxRound uint64) uint64 {
	rand.Seed(ctx.BlockTime().UnixNano())
	round := rand.Uint64() % maxRound

	store := ctx.KVStore(k.storeKey)
	it := store.Iterator(append(types.PrefixMintableMetadataIds, sdk.Uint64ToBigEndian(collId)...), nil)
	defer it.Close()

	for i := uint64(0); i < round; i++ {
		if !it.Valid() {
			return 0
		}
		it.Next()
	}

	value := sdk.BigEndianToUint64(it.Value())
	store.Delete(it.Key())
	return value
}

func (k Keeper) DeleteMintableMetadataIds(ctx sdk.Context, collId uint64) {
	store := ctx.KVStore(k.storeKey)
	it := store.Iterator(append(types.PrefixMintableMetadataIds, sdk.Uint64ToBigEndian(collId)...), nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		store.Delete(it.Key())
	}
}

func (k Keeper) ShuffleMintableMetadataIds(ctx sdk.Context, collId uint64) {
	mintableMetadataIds := k.GetMintableMetadataIds(ctx, collId)
	mintableMetadataIds = RandomList(ctx, mintableMetadataIds)
	k.DeleteMintableMetadataIds(ctx, collId)
	k.SetMintableMetadataIds(ctx, collId, mintableMetadataIds)
}

func (k Keeper) GetMintableMetadataIds(ctx sdk.Context, collId uint64) []uint64 {
	store := ctx.KVStore(k.storeKey)
	it := store.Iterator(append(types.PrefixMintableMetadataIds, sdk.Uint64ToBigEndian(collId)...), nil)
	defer it.Close()

	metadataIds := []uint64{}
	for ; it.Valid(); it.Next() {
		metadataIds = append(metadataIds, sdk.BigEndianToUint64(it.Value()))
	}

	return metadataIds
}

func RandomList(ctx sdk.Context, items []uint64) []uint64 {
	rand.Seed(ctx.BlockTime().UnixNano())
	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
	return items
}
