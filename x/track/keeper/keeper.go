package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"
	"sort"
)

// Keeper of the track store
type Keeper struct {
	bankKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}

// NewKeeper creates a content keeper
func NewKeeper(bk bank.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		bankKeeper: bk,
		storeKey:   key,
		cdc:        cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetTrack(ctx sdk.Context, id uint64) (track types.Track, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTrackKey(id))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &track)
	return track, true
}

func (k Keeper) SetTrack(ctx sdk.Context, track *types.Track) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&track)
	store.Set(types.GetTrackKey(track.TrackID), bz)
}

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

func (k Keeper) GetTracks(ctx sdk.Context) []types.Track {
	var tracks []types.Track
	k.IterateTracks(ctx, func(track types.Track) (stop bool) {
		tracks = append(tracks, track)
		return false
	})
	return tracks
}

func (k Keeper) GetTracksPaginated(ctx sdk.Context, params types.QueryTracksParams) []types.Track {
	var tracks []types.Track
	k.IterateTracks(ctx, func(track types.Track) (stop bool) {
		tracks = append(tracks, track)
		return false
	})

	sort.Slice(tracks, func(i, j int) bool {
		a, b := tracks[i], tracks[j]
		return a.TrackID > b.TrackID
	})

	page := params.Page
	if page == 0 {
		page = 1
	}
	start, end := client.Paginate(len(tracks), page, params.Limit, 100)
	if start < 0 || end < 0 {
		tracks = []types.Track{}
	} else {
		tracks = tracks[start:end]
	}

	return tracks
}

func (k Keeper) Add(ctx sdk.Context, track *types.Track) (uint64, error) {
	track.TrackID = k.autoIncrementID(ctx, types.KeyLastTrackID)
	track.Uri = k.generateTrackUri(ctx, track.TrackID)

	//content.CreatedAt = ctx.BlockHeader().Time
	k.SetTrack(ctx, track)
	k.SetCreatorTrack(ctx, track)

	return track.TrackID, nil
}

func (k Keeper) autoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(lastIDKey)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(lastIDKey, bz)
	return id
}

func (k Keeper) generateTrackUri(ctx sdk.Context, trackID uint64) string {
	return fmt.Sprintf("bitsong:track:%d", trackID)
}

func (k Keeper) SetCreatorTrack(ctx sdk.Context, track *types.Track) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&track.TrackID)
	store.Set(types.GetTrackByCreatorAddr(track.Creator, track.TrackID), bz)
}

func (k Keeper) GetCreatorTracks(ctx sdk.Context, creator sdk.AccAddress) (tracks []types.Track) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetCreatorKey(creator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var trackID uint64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &trackID)
		track, ok := k.GetTrack(ctx, trackID)
		if ok {
			tracks = append(tracks, track)
		}
	}

	return
}

/*func (k Keeper) IterateTracks(ctx sdk.Context, fn func(track types.Track) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TrackKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var track types.Track
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &track)
		if fn(track) {
			break
		}
	}
}*/

/*func allocateDaoFunds(cnt types.Content, coin sdk.Coin) types.Content {
	// allocate dao funds
	for i, rh := range cnt.RightsHolders {
		price := sdk.NewDecCoinFromCoin(coin)
		allocation := price.Amount.Quo(sdk.NewDec(100).Quo(rh.Quota))
		cnt.RightsHolders[i].Rewards = rh.Rewards.Add(sdk.NewDecCoinFromDec(btsg.BondDenom, allocation))
	}

	return cnt
}*/
