package keeper

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey sdk.StoreKey
	codec    *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, codec *codec.Codec) Keeper {
	keeper := Keeper{
		storeKey: storeKey,
		codec:    codec,
	}

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetArtist(ctx sdk.Context, id btsg.ID) (artist types.Artist, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetArtistKey(id))
	if bz == nil {
		return
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &artist)
	return artist, true
}

func (k Keeper) SetArtist(ctx sdk.Context, artist types.Artist) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(&artist)
	store.Set(types.GetArtistKey(artist.ID), bz)
}

func (k Keeper) IterateAllArtists(ctx sdk.Context, fn func(artist types.Artist) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ArtistKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var artist types.Artist
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &artist)
		if fn(artist) {
			break
		}
	}
}

func (k Keeper) GetAllArtists(ctx sdk.Context) []types.Artist {
	var artists []types.Artist
	k.IterateAllArtists(ctx, func(artist types.Artist) (stop bool) {
		artists = append(artists, artist)
		return false
	})
	return artists
}

func (k Keeper) HasID(ctx sdk.Context, id btsg.ID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetArtistKey(id))
}

func (k Keeper) CreateArtist(ctx sdk.Context, artist types.Artist) (types.Artist, error) {
	if k.HasID(ctx, artist.ID) {
		return types.Artist{}, sdkerrors.Wrap(types.ErrArtistiCreateError, fmt.Sprintf("id %s exist", artist.ID))
	}

	k.SetArtist(ctx, artist)
	k.Logger(ctx).Info(fmt.Sprintf("Artist Created %s", artist.ID))

	return artist, nil
}
