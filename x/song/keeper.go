package song

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyDelimiter = []byte(":")

	KeyNextSongId = []byte("newSongId")
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	codespace sdk.CodespaceType
}

// NewKeeper creates new instances of the song Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) AddSong(ctx sdk.Context, song Song) {
	k.setSong(ctx, song)
	idArr := k.GetAddressSongs(ctx, song.Owner)
	idArr = append(idArr, song.SongId)
	k.setAddressSongs(ctx, song.Owner, idArr)
}

func (k Keeper) Publish(ctx sdk.Context, title string,
	owner sdk.AccAddress) (song *Song, err sdk.Error) {
	id, err := k.getNewSongId(ctx)

	if err != nil {
		return nil, err
	}

	createTime := ctx.BlockHeader().Time

	song = &Song{
		SongId: id,
		Owner: owner,
		Title: title,
		CreateTime: createTime,
	}

	k.AddSong(ctx, *song)
	return song, nil
}

// Get the next available SongId and increments it
func (k Keeper) getNewSongId(ctx sdk.Context) (id uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextSongId)
	if bz == nil {
		//return 0, sdk.NewError(k.codespace, types.CodeInvalidGenesis, "InitialSongId never set")
		return 0, ErrInvalidGenesis(k.codespace)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(id + 1)
	store.Set(KeyNextSongId, bz)
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
		return nil, sdk.NewError(k.codespace, CodeSongNotExist, fmt.Sprintf("this id is invalid : %d", id))
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
func (k Keeper) PeekCurrentSongId(ctx sdk.Context) (id uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextSongId)
	if bz == nil {
		return 0, sdk.NewError(k.codespace, CodeInvalidGenesis, "InitialSongID never set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	return id, nil
}

// Set the initial song ID
func (k Keeper) SetInitialSongId(ctx sdk.Context, id uint64) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextSongId)
	if bz != nil {
		return sdk.NewError(k.codespace, CodeInvalidGenesis, "Initial SongID already set")
	}
	bz = k.cdc.MustMarshalBinaryLengthPrefixed(id)
	store.Set(KeyNextSongId, bz)
	return nil
}

func (k Keeper) setSong(ctx sdk.Context, song Song) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(song)
	store.Set(KeySong(song.SongId), bz)
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
