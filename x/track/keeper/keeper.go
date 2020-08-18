package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"
	"sort"
)

// Keeper of the track store
type Keeper struct {
	supplyKeeper supply.Keeper
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
}

// NewKeeper creates a track keeper
func NewKeeper(supplyKeeper supply.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		supplyKeeper: supplyKeeper,
		storeKey:     key,
		cdc:          cdc,
	}

	// ensure track module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetTrack(ctx sdk.Context, trackID string) (track types.Track, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTrackKey(trackID))
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

func (k Keeper) Add(ctx sdk.Context, track *types.Track) (string, error) {
	// Check to see if track exists
	_, ok := k.GetTrack(ctx, track.TrackID)
	if ok {
		return track.TrackID, sdkerrors.Wrapf(types.ErrDuplicatedTrack, "%s", track.TrackID)
	}

	//track.TrackID = k.autoIncrementID(ctx, types.KeyLastTrackID)
	//track.Uri = k.generateTrackUri(ctx, track.TrackID)

	// TODO: add created_at
	k.SetTrack(ctx, track)
	k.SetCreatorTrack(ctx, track)

	return track.TrackID, nil
}

func (k Keeper) generateTrackUri(ctx sdk.Context, trackID string) string {
	return fmt.Sprintf("bitsong:track:%s", trackID)
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
		var trackID string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &trackID)
		track, ok := k.GetTrack(ctx, trackID)
		if ok {
			tracks = append(tracks, track)
		}
	}

	return
}

func (k Keeper) MintAndSend(ctx sdk.Context, amount sdk.Coin, recipient sdk.AccAddress) error {
	// TODO: add security checks and improve

	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{amount}); err != nil {
		return err
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, sdk.Coins{amount}); err != nil {
		return err
	}

	return nil
}
