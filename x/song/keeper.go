package song

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) GetSong(ctx sdk.Context, title string) Song {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(title)) {
		return
	}
	bz := store.Get([]byte(title))
	var song Song
	k.cdc.MustUnmarshalBinaryBare(bz, &song)
	return song
}

// SetSong - sets the value string that a name resolves to
func (k Keeper) SetSong(ctx sdk.Context, title string, song Song) {
	if song.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(title), k.cdc.MustUnmarshalBinaryBare(song))
}

func (k Keeper) SetTitle(ctx sdk.Context, title string, value string) {
	song := k.GetSong(ctx, title)
	song.Value = value
	k.SetSong(ctx, title, song)
}

func (k Keeper) HasOwner(ctx sdk.Context, title string) bool {
	return !k.GetSong(ctx, title).Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, title string) sdk.AccAddress {
	return k.GetSong(ctx, title).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, title string, owner sdk.AccAddress) {
	song := k.GetSong(ctx, title)
	song.Owner = owner
	k.SetSong(ctx, title, song)
}

func (k Keeper) GetTitlesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}