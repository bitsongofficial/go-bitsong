package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/album/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddTrack Add a track on a specific album
func (k Keeper) AddTrack(ctx sdk.Context, albumID uint64, trackID uint64, position uint64) sdk.Error {
	album, ok := k.GetAlbum(ctx, albumID)
	if !ok {
		return types.ErrUnknownAlbum(k.codespace, fmt.Sprintf("unknown albumID %d", albumID))
	}

	// TODO:
	// only status NIL ?
	if album.Status != types.StatusNil {
		return types.ErrInvalidAlbumStatus(k.codespace, fmt.Sprintf("album status must be nil"))
	}

	// TODO:
	// check if track exist
	// check if track status is nil
	// check track is on duplicate
	// check position is allowed

	track := types.NewTrack(albumID, trackID, position)
	k.setTrack(ctx, albumID, track)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddTrackAlbum,
			sdk.NewAttribute(types.AttributeKeyAlbumID, fmt.Sprintf("%d", albumID)),
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", track.TrackID)),
		),
	)

	return nil
}

func (k Keeper) setTrack(ctx sdk.Context, albumID uint64, track types.Track) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(track)
	store.Set(types.TrackKey(albumID, track.TrackID), bz)
}

// GetTracksIterator gets all the tracks on a specific album as an sdk.Iterator
func (k Keeper) GetTracksIterator(ctx sdk.Context, albumID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.TracksKey(albumID))
}

// IterateTracks iterates over the all the albums tracks and performs a callback function
func (k Keeper) IterateTracks(ctx sdk.Context, albumID uint64, cb func(track types.Track) (stop bool)) {
	iterator := k.GetTracksIterator(ctx, albumID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var track types.Track
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &track)

		if cb(track) {
			break
		}
	}
}

// GetTracks returns all the tracks from a specific album
func (k Keeper) GetTracks(ctx sdk.Context, albumID uint64) (tracks types.Tracks) {
	k.IterateTracks(ctx, albumID, func(track types.Track) bool {
		tracks = append(tracks, track)
		return false
	})
	return
}

// GetTrack gets the track from  specific trackID on a specific album
func (k Keeper) GetTrack(ctx sdk.Context, albumID uint64, trackID uint64) (track types.Track, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TrackKey(albumID, trackID))
	if bz == nil {
		return track, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &track)
	return track, true
}
