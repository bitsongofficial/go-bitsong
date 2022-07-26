package keeper

import (
	"bytes"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

/*func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

func (k Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coins {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	return k.bankKeeper.GetAllBalances(ctx, moduleAccAddr)
}*/

func (k Keeper) SetLastMerkleDropId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastMerkledropIDKey(), sdk.Uint64ToBigEndian(id))
}

func (k Keeper) GetLastMerkleDropId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastMerkledropIDKey())
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetMerkleDrop(ctx sdk.Context, merkledrop types.Merkledrop) error {
	bz := k.cdc.MustMarshal(&merkledrop)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MerkledropKey(merkledrop.Id), bz)

	owner, err := sdk.AccAddressFromBech32(merkledrop.Owner)
	if err != nil {
		return err
	}

	// set key by owner
	store.Set(types.MerkledropOwnerKey(merkledrop.Id, owner), sdk.Uint64ToBigEndian(merkledrop.Id))

	// set key by end-height
	store.Set(types.MerkledropEndHeightAndIDKey(merkledrop.EndHeight, merkledrop.Id), []byte{0x01})

	return nil
}

func (k Keeper) IsClaimed(ctx sdk.Context, mdId, index uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ClaimedMerkledropIndexKey(mdId, index))
}

func (k Keeper) SetClaimed(ctx sdk.Context, mdId, index uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ClaimedMerkledropIndexKey(mdId, index), []byte{0x01})
}

func (k Keeper) GetAllMerkleDrops(ctx sdk.Context) []types.Merkledrop {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.PrefixMerkleDrop)
	defer it.Close()

	var allMerkleDrops []types.Merkledrop
	for ; it.Valid(); it.Next() {
		var merkledrop types.Merkledrop
		k.cdc.MustUnmarshal(it.Value(), &merkledrop)

		allMerkleDrops = append(allMerkleDrops, merkledrop)
	}

	return allMerkleDrops
}

func (k Keeper) getMerkleDropById(ctx sdk.Context, id uint64) (types.Merkledrop, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MerkledropKey(id))
	if bz == nil {
		return types.Merkledrop{}, sdkerrors.Wrapf(types.ErrMerkledropNotExist, "merkledrop: %d does not exist", id)
	}
	merkledrop := types.Merkledrop{}
	k.cdc.MustUnmarshal(bz, &merkledrop)
	return merkledrop, nil
}

func (k Keeper) getMerkleDropsByOwner(ctx sdk.Context, owner sdk.AccAddress) []types.Merkledrop {
	store := ctx.KVStore(k.storeKey)

	var merkledrops []types.Merkledrop
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixMerkleDropByOwner, owner...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		merkledrop, err := k.getMerkleDropById(ctx, id)
		if err != nil {
			panic(err)
		}

		merkledrops = append(merkledrops, merkledrop)
	}
	return merkledrops
}

func (k Keeper) GetMerkleDropsIDByEndHeight(ctx sdk.Context, endHeight int64) []uint64 {
	var mdIDs []uint64
	k.iterateMerkledropIDByEndHeight(ctx, endHeight, func(mdID uint64) (stop bool) {
		mdIDs = append(mdIDs, mdID)
		return false
	})

	return mdIDs
}

func (k Keeper) iterateMerkledropIDByEndHeight(ctx sdk.Context, endHeight int64, cb func(mdID uint64) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.MerkledropEndHeightKey(endHeight)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		mdID := sdk.BigEndianToUint64(bytes.TrimPrefix(iterator.Key(), prefix))

		if cb(mdID) {
			break
		}
	}
}

func (k Keeper) iterateIndexByMerkledropID(ctx sdk.Context, mdId uint64, cb func(index uint64) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.ClaimedMerkledropKey(mdId)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		index := sdk.BigEndianToUint64(bytes.TrimPrefix(iterator.Key(), prefix))

		if cb(index) {
			break
		}
	}
}

func (k Keeper) GetAllIndexesByMerkledropID(ctx sdk.Context, id uint64) []uint64 {
	var indexes []uint64
	k.iterateIndexByMerkledropID(ctx, id, func(index uint64) (stop bool) {
		indexes = append(indexes, index)
		return false
	})

	return indexes
}

func (k Keeper) deleteAllIndexesByMerkledropID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	indexes := k.GetAllIndexesByMerkledropID(ctx, id)

	for _, index := range indexes {
		store.Delete(types.ClaimedMerkledropIndexKey(id, index))
	}
}

func (k Keeper) DeleteMerkledropByID(ctx sdk.Context, id uint64) error {
	store := ctx.KVStore(k.storeKey)

	merkledrop, err := k.getMerkleDropById(ctx, id)
	if err != nil {
		return err
	}

	// delete owner store
	owner, err := sdk.AccAddressFromBech32(merkledrop.Owner)
	if err != nil {
		return err
	}
	store.Delete(types.MerkledropOwnerKey(merkledrop.Id, owner))

	// delete all indexes
	k.deleteAllIndexesByMerkledropID(ctx, id)

	// delete end-height key
	store.Delete(types.MerkledropEndHeightAndIDKey(merkledrop.EndHeight, merkledrop.Id))

	// delete merkledrop
	store.Delete(types.MerkledropKey(id))

	return nil
}

func (k Keeper) GetAllIndexes(ctx sdk.Context) []*types.Indexes {
	var mdIndexes []*types.Indexes
	merkledrops := k.GetAllMerkleDrops(ctx)

	for _, md := range merkledrops {
		var indexes []uint64

		k.iterateIndexByMerkledropID(ctx, md.Id, func(index uint64) (stop bool) {
			indexes = append(indexes, index)
			return false
		})

		mdIndexes = append(mdIndexes, &types.Indexes{
			MerkledropId: md.Id,
			Index:        indexes,
		})
	}
	return mdIndexes
}
