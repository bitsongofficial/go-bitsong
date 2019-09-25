package song

import (
	"fmt"
	"github.com/BitSongOfficial/go-bitsong/x/song/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyDelimiter = []byte(":")

	KeyNextSongID = []byte("newSongID")
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
	paramSpace params.Subspace
	sk         staking.Keeper
}

// NewKeeper creates new instances of the song Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, paramSpace params.Subspace, sk staking.Keeper) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		sk:         sk,
	}
}

// AddSong keeper
func (k Keeper) AddSong(ctx sdk.Context, song Song) {
	k.setSong(ctx, song)
	idArr := k.GetAddressSongs(ctx, song.Owner)
	idArr = append(idArr, song.SongID)
	k.setAddressSongs(ctx, song.Owner, idArr)
}

func (k Keeper) Play(ctx sdk.Context, songId uint64, accAddr sdk.AccAddress) sdk.Error {
	// Get user power
	//delegations := k.sk.GetAllDelegatorDelegations(ctx, accAddr)

	return nil
}

func (k Keeper) GetPlay(ctx sdk.Context, songId uint64, accAddr sdk.AccAddress) (play *types.Play, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyPlay(accAddr, songId))
	if bz == nil {
		return nil, sdk.NewError(DefaultCodespace, CodePlayNotExist, fmt.Sprintf("this play is invalid"))
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &play)
	return play, nil
}

func (k Keeper) SetPlay(ctx sdk.Context, songId uint64, accAddr sdk.AccAddress) (play *types.Play) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyPlay(accAddr, songId))
	if bz == nil {
		play = &types.Play{
			AccountAddress: accAddr,
			SongId:         songId,
			Shares:         sdk.Dec{},
			Streams:        sdk.NewInt(1),
		}
	}

	play.Streams = play.Streams.Add(sdk.NewInt(1))

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &play)
	return play
}

func KeyPlay(accAddr sdk.AccAddress, songId uint64) []byte {
	return []byte(fmt.Sprintf("play:%s-%d", accAddr, songId))
}

// Publish keeper
func (k Keeper) Publish(ctx sdk.Context, title string,
	owner sdk.AccAddress, content string,
	redistributionSplitRate sdk.Dec) (song *Song, err sdk.Error) {
	id, err := k.getNewSongID(ctx)

	if err != nil {
		return nil, err
	}

	createTime := ctx.BlockHeader().Time
	totalReward := sdk.NewInt(0)

	song = &Song{
		SongID:                  id,
		Owner:                   owner,
		Title:                   title,
		Content:                 content,
		TotalReward:             totalReward,
		RedistributionSplitRate: redistributionSplitRate,
		CreateTime:              createTime,
	}

	k.AddSong(ctx, *song)
	return song, nil
}

// Get the next available SongID and increments it
func (k Keeper) getNewSongID(ctx sdk.Context) (id uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextSongID)
	if bz == nil {
		//return 0, sdk.NewError(k.codespace, types.CodeInvalidGenesis, "InitialSongID never set")
		return 0, ErrInvalidGenesis(DefaultCodespace)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(id + 1)
	store.Set(KeyNextSongID, bz)
	return id, nil
}

func (k Keeper) GetSongsByAddr(ctx sdk.Context, addr sdk.AccAddress) (songs Songs, err sdk.Error) {
	idArr := k.GetAddressSongs(ctx, addr)
	for _, id := range idArr {
		song, err := k.GetSong(ctx, id)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// Store level
func (k Keeper) GetSong(ctx sdk.Context, id uint64) (*Song, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeySong(id))
	if bz == nil {
		return nil, sdk.NewError(DefaultCodespace, CodeSongNotExist, fmt.Sprintf("this id is invalid : %d", id))
	}
	var song Song
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &song)
	return &song, nil
}

// Key for getting a specific song from the store
func KeySong(id uint64) []byte {
	return []byte(fmt.Sprintf("songs:%d", id))
}

// Key for getting all songs of a owner from the store
func KeyAddressSongs(addr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("address:%d", addr))
}

// Peeks the next available id without incrementing it
func (k Keeper) PeekCurrentSongID(ctx sdk.Context) (id uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextSongID)
	if bz == nil {
		return 0, sdk.NewError(DefaultCodespace, CodeInvalidGenesis, "InitialSongID never set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	return id, nil
}

// Set the initial song ID
func (k Keeper) SetInitialSongID(ctx sdk.Context, id uint64) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextSongID)
	if bz != nil {
		return sdk.NewError(DefaultCodespace, CodeInvalidGenesis, "Initial SongID already set")
	}
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(id)
	store.Set(KeyNextSongID, bz)
	return nil
}

func (k Keeper) setSong(ctx sdk.Context, song Song) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(song)
	store.Set(KeySong(song.SongID), bz)
}

func (k Keeper) SetSong(ctx sdk.Context, song Song) {
	k.setSong(ctx, song)
}

func (k Keeper) SetAddressSong(ctx sdk.Context, addr sdk.AccAddress) (idArr []uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyAddressSongs(addr))
	if bz == nil {
		return []uint64{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &idArr)
	return idArr
}

func (k Keeper) setAddressSongs(ctx sdk.Context, addr sdk.AccAddress, idArr []uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(idArr)
	store.Set(KeyAddressSongs(addr), bz)
}

func (k Keeper) GetAddressSongs(ctx sdk.Context, addr sdk.AccAddress) (idArr []uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyAddressSongs(addr))
	if bz == nil {
		return []uint64{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &idArr)
	return idArr
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
