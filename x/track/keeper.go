package track

import (
	"fmt"
	"github.com/BitSongOfficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"

	//"github.com/cosmos/cosmos-sdk/x/staking/exported"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyDelimiter                = []byte(":")
	KeyNextTrackID              = []byte("newTrackID")
	PlayPrefix                  = []byte("play")
	AccountCurrentRewardsPrefix = []byte("accCurrReward")
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSpace    params.Subspace
	stakingKeeper staking.Keeper
	supplyKeeper  supply.Keeper
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, paramSpace params.Subspace, stkingKeeper staking.Keeper, supplyKeeper supply.Keeper) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSpace:    paramSpace.WithKeyTable(types.ParamKeyTable()),
		stakingKeeper: stkingKeeper,
		supplyKeeper:  supplyKeeper,
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

/*func KeyPlay(accAddr sdk.AccAddress, trackID uint64) []byte {
	return append([]byte(PlayPrefix), accAddr..., trackID...)
}*/

func GetPlayKey(accAddr sdk.AccAddress, trackID uint64) []byte {
	return append(GetPlaysKey(trackID), accAddr.Bytes()...)
}

func GetPlaysKey(trackID uint64) []byte {
	return append(PlayPrefix, make([]byte, trackID)...)
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

func (k Keeper) GetAccPower(ctx sdk.Context, address sdk.AccAddress) sdk.Dec {
	power := sdk.ZeroDec()

	k.stakingKeeper.IterateDelegations(
		ctx, address,
		func(_ int64, del exported.DelegationI) (stop bool) {
			power = power.Add(del.GetShares())
			return false
		},
	)

	return power
}

func (k Keeper) GetPlays(ctx sdk.Context, trackId uint64) (types.Play, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetPlaysKey(trackId))
	if bz == nil {
		return types.Play{}, false
	}
	var play types.Play
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &play)
	return play, true
}

func (k Keeper) GetAccPlay(ctx sdk.Context, accAddr sdk.AccAddress, trackId uint64) (types.Play, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetPlayKey(accAddr, trackId))
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
	store.Set(GetPlayKey(play.AccAddress, play.TrackId), bz)

	return nil
}

func (k Keeper) SetPlay(ctx sdk.Context, play types.Play) sdk.Error {
	return k.setPlay(ctx, play)
}

func (k Keeper) DeletePlay(ctx sdk.Context, play types.Play) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetPlayKey(play.AccAddress, play.TrackId))
}

func (k Keeper) PlayTrack(ctx sdk.Context, accAddr sdk.AccAddress, trackID uint64) (types.Play, bool) {
	play, ok := k.GetAccPlay(ctx, accAddr, trackID)

	if !ok {
		play = types.Play{
			AccAddress: accAddr,
			TrackId:    trackID,
			Shares:     k.GetAccPower(ctx, accAddr),
			Streams:    sdk.NewInt(1),
			CreateTime: ctx.BlockHeader().Time,
		}

		k.setPlay(ctx, play)
	} else {
		play.Streams = play.Streams.Add(sdk.NewInt(1))
		k.setPlay(ctx, play)
	}

	return play, true
}

func (k Keeper) IterateAllPlays(ctx sdk.Context, cb func(play types.Play) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, PlayPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		play, err := UnmarshalPlay(k.cdc, iterator.Value())
		if err != nil {
			panic(err)
		}

		if cb(play) {
			break
		}
	}
}

func UnmarshalPlay(cdc *codec.Codec, value []byte) (play types.Play, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &play)
	return play, err
}

func (k Keeper) GetAllPlays(ctx sdk.Context) (plays []types.Play) {
	k.IterateAllPlays(ctx, func(play types.Play) bool {
		plays = append(plays, play)
		return false
	})
	return plays
}

func (k Keeper) GetPlayTax(ctx sdk.Context) sdk.Dec {
	var percent sdk.Dec
	k.paramSpace.Get(ctx, types.KeyPlayTax, &percent)
	return percent
}

func (k Keeper) GetFeePlayPool(ctx sdk.Context) (playPool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.PlayPoolKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &playPool)
	return
}

func (k Keeper) SetFeePlayPool(ctx sdk.Context, playPool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(playPool)
	store.Set(types.PlayPoolKey, b)
}

// gets the key for a account's current rewards
func GetAccountCurrentRewardsKey(accAddr sdk.AccAddress) []byte {
	return append(AccountCurrentRewardsPrefix, accAddr.Bytes()...)
}

// gets the address
func GetAccountCurrentRewardsAddress(key []byte) (accAddr sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return addr
}

// get current rewards for an account
func (k Keeper) GetAccountCurrentRewards(ctx sdk.Context, acc sdk.AccAddress) (rewards types.AccountCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetAccountCurrentRewardsKey(acc))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &rewards)
	return
}

// set current rewards for an account
func (k Keeper) SetAccountCurrentRewards(ctx sdk.Context, acc sdk.AccAddress, rewards types.AccountCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(rewards)
	store.Set(GetAccountCurrentRewardsKey(acc), b)
}

// delete current rewards for an account
func (k Keeper) DeleteAccountCurrentRewards(ctx sdk.Context, acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetAccountCurrentRewardsKey(acc))
}

// iterate over current rewards
func (k Keeper) IterateAccountCurrentRewards(ctx sdk.Context, handler func(acc sdk.AccAddress, rewards types.AccountCurrentRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, AccountCurrentRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.AccountCurrentRewards
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &rewards)
		addr := GetAccountCurrentRewardsAddress(iter.Key())
		if handler(addr, rewards) {
			break
		}
	}
}

func (k Keeper) AllocateTokensToAccount(ctx sdk.Context, acc sdk.AccAddress, tokens sdk.DecCoins) {
	// update current rewards
	currentRewards := k.GetAccountCurrentRewards(ctx, acc)
	currentRewards.Rewards = currentRewards.Rewards.Add(tokens)
	k.SetAccountCurrentRewards(ctx, acc, currentRewards)
}
