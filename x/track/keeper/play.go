package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
)

// GetAccPower, get account power based on staking
func (keeper Keeper) GetAccPower(ctx sdk.Context, address sdk.AccAddress) sdk.Dec {
	power := sdk.ZeroDec()

	keeper.stakingKeeper.IterateDelegations(
		ctx, address,
		func(_ int64, del exported.DelegationI) (stop bool) {
			power = power.Add(del.GetShares())
			return false
		},
	)

	return power
}

// Play Add a play on a specific track
func (keeper Keeper) Play(ctx sdk.Context, trackID uint64, accAddr sdk.AccAddress) sdk.Error {
	track, ok := keeper.GetTrack(ctx, trackID)
	if !ok {
		return types.ErrUnknownTrack(keeper.codespace, fmt.Sprintf("unknown trackID %d", trackID))
	}

	// TODO:
	// only status VERIFIED ?
	if track.Status != types.StatusVerified {
		return types.ErrInvalidTrackStatus(keeper.codespace, fmt.Sprintf("track status must be verified"))
	}

	// TODO:
	// improve checks

	play, ok := keeper.GetPlay(ctx, trackID, accAddr)

	createdAt := ctx.BlockHeader().Time
	streams := uint64(1)
	shares := keeper.GetAccPower(ctx, accAddr)
	if !shares.IsPositive() {
		// TODO: change error
		return types.ErrInvalidTrackStatus(keeper.codespace, fmt.Sprintf("user share must be positive"))
	}

	if !ok {
		play = types.NewPlay(
			trackID,
			accAddr,
			shares,
			streams,
			createdAt,
		)

		keeper.IncrementShare(ctx, trackID, shares)
	} else {
		play.Streams = play.Streams + streams
	}

	keeper.setPlay(ctx, trackID, play)

	// TODO:
	// improve increment
	track.TotalPlays = track.TotalPlays + 1
	keeper.SetTrack(ctx, track)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePlayTrack,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", track.TrackID)),
			sdk.NewAttribute(types.AttributeKeyPlayAccAddr, fmt.Sprintf("%s", play.AccAddr.String())),
		),
	)

	return nil
}

func (keeper Keeper) setPlay(ctx sdk.Context, trackID uint64, play types.Play) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(play)
	store.Set(types.PlayKey(trackID, play.AccAddr), bz)
}

func (keeper Keeper) GetPlaysIterator(ctx sdk.Context, trackID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.PlaysKey(trackID))
}

// IterateAllPlays iterates over the all the stored plays and performs a callback function
func (keeper Keeper) IterateAllPlays(ctx sdk.Context, cb func(play types.Play) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PlaysKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var play types.Play
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &play)

		if cb(play) {
			break
		}
	}
}

func (keeper Keeper) IteratePlays(ctx sdk.Context, trackID uint64, cb func(play types.Play) (stop bool)) {
	iterator := keeper.GetPlaysIterator(ctx, trackID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var play types.Play
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &play)

		if cb(play) {
			break
		}
	}
}

func (keeper Keeper) GetPlays(ctx sdk.Context, trackID uint64) (plays types.Plays) {
	keeper.IteratePlays(ctx, trackID, func(play types.Play) bool {
		plays = append(plays, play)
		return false
	})
	return
}

func (keeper Keeper) GetPlay(ctx sdk.Context, trackID uint64, accAddr sdk.AccAddress) (play types.Play, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.PlayKey(trackID, accAddr))
	if bz == nil {
		return play, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &play)
	return play, true
}
