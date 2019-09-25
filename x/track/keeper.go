package track

import (
	"fmt"
	"github.com/BitSongOfficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	//"github.com/cosmos/cosmos-sdk/x/staking/exported"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyDelimiter = []byte(":")

	KeyNextTrackID = []byte("newTrackID")
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace params.Subspace
	sk         staking.Keeper
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, paramSpace params.Subspace, sk staking.Keeper) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		sk:         sk,
	}
}

// Set the initial track ID
func (k Keeper) SetInitialTrackID(ctx sdk.Context, id uint64) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextTrackID)
	if bz != nil {
		return ErrInitialTrackIDAlreadySet(DefaultCodespace)
	}
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(id)
	store.Set(KeyNextTrackID, bz)
	return nil
}

// Get the next available TrackID and increments it
func (k Keeper) getNewTrackID(ctx sdk.Context) (id uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextTrackID)
	if bz == nil {
		return 0, ErrInvalidGenesis(DefaultCodespace)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(id + 1)
	store.Set(KeyNextTrackID, bz)
	return id, nil
}

// Peeks the next available id without incrementing it
func (k Keeper) PeekCurrentTrackID(ctx sdk.Context) (id uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextTrackID)
	if bz == nil {
		return 0, ErrInvalidGenesis(DefaultCodespace)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	return id, nil
}

// Key for getting a specific track from the store
func KeyTrack(id uint64) []byte {
	return []byte(fmt.Sprintf("track:%d", id))
}

func KeyPlay(accAddr sdk.AccAddress, trackID uint64) []byte {
	return []byte(fmt.Sprintf("play:%s-%d", accAddr, trackID))
}

func (k Keeper) GetTrack(ctx sdk.Context, trackId uint64) (types.Track, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyTrack(trackId))
	if bz == nil {
		return types.Track{}, false
	}
	var track types.Track
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &track)
	return track, true
}

func (k Keeper) setTrack(ctx sdk.Context, track types.Track) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(track)
	store.Set(KeyTrack(track.TrackID), bz)

	return nil
}

func (k Keeper) SetTrack(ctx sdk.Context, track types.Track) sdk.Error {
	return k.setTrack(ctx, track)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) PublishTrack(ctx sdk.Context, title string, owner sdk.AccAddress, content string, redistributionSplitRate sdk.Dec) (track *Track, err sdk.Error) {
	id, err := k.getNewTrackID(ctx)
	if err != nil {
		return nil, err
	}

	createTime := ctx.BlockHeader().Time
	totalReward := sdk.NewInt(0)

	track = &Track{
		TrackID:                 id,
		Owner:                   owner,
		Title:                   title,
		Content:                 content,
		TotalReward:             totalReward,
		RedistributionSplitRate: redistributionSplitRate,
		CreateTime:              createTime,
	}

	k.setTrack(ctx, *track)
	return track, nil
}

func (k Keeper) GetUserPower(ctx sdk.Context, address sdk.AccAddress) sdk.Dec {
	power := sdk.ZeroDec()

	k.sk.IterateDelegations(
		ctx, address,
		func(_ int64, del exported.DelegationI) (stop bool) {
			power = power.Add(del.GetShares())
			return false
		},
	)

	return power
}

func (k Keeper) GetPlay(ctx sdk.Context, accAddr sdk.AccAddress, trackId uint64) (types.Play, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyPlay(accAddr, trackId))
	if bz == nil {
		return types.Play{}, false
	}
	var play types.Play
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &play)
	return play, true
}

func (k Keeper) setPlay(ctx sdk.Context, play types.Play) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(play)
	store.Set(KeyPlay(play.AccAddress, play.TrackId), bz)

	return nil
}
