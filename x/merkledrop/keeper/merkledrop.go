package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

func (k Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coins {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	return k.bankKeeper.GetAllBalances(ctx, moduleAccAddr)
}

func (k Keeper) CreateModuleAccount(ctx sdk.Context) {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
}

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

func (k Keeper) SetMerkleDrop(ctx sdk.Context, merkledrop types.Merkledrop) {
	bz := k.cdc.MustMarshal(&merkledrop)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MerkledropKey(merkledrop.Id), bz)

	owner, err := sdk.AccAddressFromBech32(merkledrop.Owner)
	if err != nil {
		panic(err)
	}
	store.Set(types.MerkledropOwnerKey(merkledrop.Id, owner), sdk.Uint64ToBigEndian(merkledrop.Id))
}

func (k Keeper) IsClaimed(ctx sdk.Context, mdId, index uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ClaimedMerkledropKey(mdId, index))
}

func (k Keeper) SetClaimed(ctx sdk.Context, mdId, index uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ClaimedMerkledropKey(mdId, index), []byte{0x01})
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

func (k Keeper) GetMerkleDropById(ctx sdk.Context, id uint64) (types.Merkledrop, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MerkledropKey(id))
	if bz == nil {
		return types.Merkledrop{}, sdkerrors.Wrapf(types.ErrMerkleDropNotExist, "merkledrop: %d does not exist", id)
	}
	merkledrop := types.Merkledrop{}
	k.cdc.MustUnmarshal(bz, &merkledrop)
	return merkledrop, nil
}

func (k Keeper) GetMerkleDropsByOwner(ctx sdk.Context, owner sdk.AccAddress) []types.Merkledrop {
	store := ctx.KVStore(k.storeKey)

	var merkledrops []types.Merkledrop
	it := sdk.KVStorePrefixIterator(store, append(types.PrefixMerkleDropByOwner, owner...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		id := sdk.BigEndianToUint64(it.Value())
		merkledrop, err := k.GetMerkleDropById(ctx, id)
		if err != nil {
			panic(err)
		}

		merkledrops = append(merkledrops, merkledrop)
	}
	return merkledrops
}