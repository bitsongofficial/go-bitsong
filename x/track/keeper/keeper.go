package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the track store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper creates a track keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) autoIncrementID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastTrackID)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(types.KeyLastTrackID, bz)
	return id
}

func (k Keeper) GetLastTrackID(ctx sdk.Context) (lastTrackID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastTrackID)
	if bz == nil {
		return 0, fmt.Errorf("initial track ID hasn't been set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &lastTrackID)
	return lastTrackID, nil
}

func (k Keeper) SetLastTrackID(ctx sdk.Context, lastTrackID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(lastTrackID)
	store.Set(types.KeyLastTrackID, bz)
}

func generateTrackAddress(id uint64) crypto.Address {
	addr := make([]byte, 20)
	addr[0] = 'T' // TrackAddress prefix
	binary.PutUvarint(addr[1:], id)
	return crypto.AddressHash(addr)
}

// GetTrack get Track from store by TrackID
func (k Keeper) GetTrack(ctx sdk.Context, addr crypto.Address) (track types.Track, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTrackKey(addr))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &track)
	return track, true
}

func (k Keeper) SetTrack(ctx sdk.Context, track types.Track) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(track)
	store.Set(types.GetTrackKey(track.Address), bz)
}

// IterateTracks iterates through the tracks set and performs the provided function
func (k Keeper) IterateTracks(ctx sdk.Context, fn func(track types.Track) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TrackKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var track types.Track
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &track)
		if fn(track) {
			break
		}
	}
}

func (k Keeper) GetTracks(ctx sdk.Context) types.Tracks {
	var tracks types.Tracks
	k.IterateTracks(ctx, func(track types.Track) (stop bool) {
		tracks = append(tracks, track)
		return false
	})
	return tracks
}

func (k Keeper) Create(ctx sdk.Context, track types.Track) crypto.Address {
	trackID := k.autoIncrementID(ctx)

	track.Address = generateTrackAddress(trackID)
	track.SubmitTime = ctx.BlockHeader().Time

	k.SetTrack(ctx, track)

	return track.Address
}
