package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Play Add a play on a specific track
func (k Keeper) Play(ctx sdk.Context, trackID uint64, accAddr sdk.AccAddress) sdk.Error {
	track, ok := k.GetTrack(ctx, trackID)
	if !ok {
		return types.ErrUnknownTrack(k.codespace, fmt.Sprintf("unknown trackID %d", trackID))
	}

	// TODO:
	// only status VERIFIED ?
	if track.Status != types.StatusVerified {
		return types.ErrInvalidTrackStatus(k.codespace, fmt.Sprintf("track status must be verified"))
	}

	// TODO:
	// improve checks

	createdAt := sdk.NewInt(ctx.BlockHeight())

	play := types.NewPlay(trackID, accAddr, createdAt)
	k.setPlay(ctx, trackID, play)

	// TODO:
	// improve increment
	track.TotalPlays = track.TotalPlays + 1
	k.SetTrack(ctx, track)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePlayTrack,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", track.TrackID)),
			sdk.NewAttribute(types.AttributeKeyPlayAccAddr, fmt.Sprintf("%s", play.AccAddr.String())),
		),
	)

	return nil
}

func (k Keeper) setPlay(ctx sdk.Context, trackID uint64, play types.Play) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(play)
	store.Set(types.PlayKey(trackID, play.AccAddr), bz)
}

func (k Keeper) GetPlaysIterator(ctx sdk.Context, trackID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.PlaysKey(trackID))
}

func (k Keeper) IteratePlays(ctx sdk.Context, trackID uint64, cb func(play types.Play) (stop bool)) {
	iterator := k.GetPlaysIterator(ctx, trackID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var play types.Play
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &play)

		if cb(play) {
			break
		}
	}
}

func (k Keeper) GetPlays(ctx sdk.Context, trackID uint64) (plays types.Plays) {
	k.IteratePlays(ctx, trackID, func(play types.Play) bool {
		plays = append(plays, play)
		return false
	})
	return
}

func (k Keeper) GetPlay(ctx sdk.Context, trackID uint64, accAddr sdk.AccAddress) (play types.Play, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PlayKey(trackID, accAddr))
	if bz == nil {
		return play, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &play)
	return play, true
}
