package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetShare gets the share of a specific entity on a specific track
func (k Keeper) GetShare(ctx sdk.Context, trackID string, entityAddr sdk.AccAddress) (share types.Share, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ShareKey(trackID, entityAddr))
	if bz == nil {
		return share, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &share)
	return share, true
}

// SetShare sets a Share to the track store
func (k Keeper) SetShare(ctx sdk.Context, share types.Share) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(share)
	store.Set(types.ShareKey(share.TrackID, share.Entity), bz)
}

// GetShares returns all the shares from a track
func (k Keeper) GetShares(ctx sdk.Context, trackID string) (shares []types.Share) {
	k.IterateShares(ctx, trackID, func(share types.Share) bool {
		shares = append(shares, share)
		return false
	})
	return
}

// IterateShares iterates over the all the tracks shares and performs a callback function
func (k Keeper) IterateShares(ctx sdk.Context, trackID string, cb func(share types.Share) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SharesKey(trackID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var share types.Share
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &share)

		if cb(share) {
			break
		}
	}
}

func (k Keeper) AddShare(ctx sdk.Context, trackID string, entityAddr sdk.AccAddress, shareAmount sdk.Coin) error {
	// Check to see if track exists
	track, ok := k.GetTrack(ctx, trackID)
	if !ok {
		return sdkerrors.Wrapf(types.ErrUnknownTrack, "%s", trackID)
	}

	// TODO: add security checks and validation

	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, entityAddr, types.ModuleName, sdk.Coins{shareAmount}); err != nil {
		return err
	}

	// Update track
	track.TotalShares = track.TotalShares.Add(shareAmount)
	k.SetTrack(ctx, &track)

	// Add or update share
	share, found := k.GetShare(ctx, trackID, entityAddr)
	if found {
		share.Shares = share.Shares.Add(shareAmount)
	} else {
		share = types.Share{
			TrackID: trackID,
			Entity:  entityAddr,
			Shares:  shareAmount,
		}
	}

	k.SetShare(ctx, share)
	return nil
}

func (k Keeper) RemoveShare(ctx sdk.Context, trackID string, entityAddr sdk.AccAddress, shareAmount sdk.Coin) error {
	// TODO: add security checks and validation

	// Check to see if track exists
	track, ok := k.GetTrack(ctx, trackID)
	if !ok {
		return sdkerrors.Wrapf(types.ErrUnknownTrack, "%s", trackID)
	}

	// Check to see if share exists
	share, ok := k.GetShare(ctx, trackID, entityAddr)
	if !ok {
		return sdkerrors.Wrapf(types.ErrUnknownShare, "%s", trackID)
	}

	// Check to see if amount exists
	if shareAmount.IsLT(share.Shares) {
		return sdkerrors.Wrapf(types.ErrInvalidAmount, "%s", shareAmount)
	}

	// Send share to account
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, entityAddr, sdk.Coins{shareAmount}); err != nil {
		return err
	}

	// Update track
	track.TotalShares = track.TotalShares.Sub(shareAmount)
	k.SetTrack(ctx, &track)

	// Update share
	share.Shares = share.Shares.Sub(shareAmount)
	k.SetShare(ctx, share)

	return nil
}
