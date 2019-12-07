package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) setShare(ctx sdk.Context, trackID uint64, share types.Share) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(share)
	store.Set(types.ShareKey(trackID), bz)
}

func (keeper Keeper) GetShare(ctx sdk.Context, trackID uint64) (share types.Share, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ShareKey(trackID))
	if bz == nil {
		return share, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &share)
	return share, true
}

func (keeper Keeper) IncrementShare(ctx sdk.Context, trackID uint64, amt sdk.Dec) sdk.Error {
	// TODO:
	// improve checks

	share, ok := keeper.GetShare(ctx, trackID)
	if !ok {
		share = types.NewShare(trackID)
	}

	share.TotalShare = share.TotalShare.Add(amt)
	keeper.setShare(ctx, trackID, share)

	return nil
}

func (keeper Keeper) IterateAllShares(ctx sdk.Context, cb func(share types.Share) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SharesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var share types.Share
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &share)

		if cb(share) {
			break
		}
	}
}

func (keeper Keeper) GetAllShares(ctx sdk.Context) (shares types.Shares) {
	keeper.IterateAllShares(ctx, func(share types.Share) bool {
		shares = append(shares, share)
		return false
	})
	return
}

func (keeper Keeper) DeleteAllShares(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateAllShares(ctx, func(share types.Share) bool {
		store.Delete(types.ShareKey(share.TrackID))
		return false
	})
}
