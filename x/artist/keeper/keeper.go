package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey  sdk.StoreKey      // The (unexposed) keys used to access the stores from the Context.
	cdc       *codec.Codec      // The codec for binary encoding/decoding.
	codespace sdk.CodespaceType // Reserved codespace
}

// NewKeeper returns an artist keeper.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

/****************************************
 * Artists
 ****************************************/

// Set the artist ID
func (keeper Keeper) SetArtistID(ctx sdk.Context, artistID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(artistID)
	store.Set(types.ArtistIDKey, bz)
}

// GetArtistID gets the highest artist ID
func (keeper Keeper) GetArtistID(ctx sdk.Context) (artistID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ArtistIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis(keeper.codespace, "initial artist ID hasn't been set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &artistID)
	return artistID, nil
}

// SetArtist set an artist to store
func (keeper Keeper) SetArtist(ctx sdk.Context, artist types.Artist) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(artist)
	store.Set(types.ArtistKey(artist.ArtistID), bz)
}

// GetArtist get Artist from store by ArtistID
func (keeper Keeper) GetArtist(ctx sdk.Context, artistID uint64) (artist types.Artist, ok bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ArtistKey(artistID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &artist)
	return artist, true
}

// GetArtistsFiltered get Artists from store by ArtistID
// status will filter artists by status
// numLatest will fetch a specified number of the most recent artists, or 0 for all artists
func (keeper Keeper) GetArtistsFiltered(ctx sdk.Context, ownerAddr sdk.AccAddress, status types.ArtistStatus, numLatest uint64) []types.Artist {

	maxArtistID, err := keeper.GetArtistID(ctx)
	if err != nil {
		return []types.Artist{}
	}

	var matchingArtists []types.Artist

	if numLatest == 0 {
		numLatest = maxArtistID
	}

	for artistID := maxArtistID - numLatest; artistID < maxArtistID; artistID++ {
		artist, ok := keeper.GetArtist(ctx, artistID)
		if !ok {
			continue
		}

		if artist.Status.Valid() && artist.Status != status {
			continue
		}

		if ownerAddr != nil && len(ownerAddr) != 0 && artist.Owner.String() != ownerAddr.String() {
			continue
		}

		matchingArtists = append(matchingArtists, artist)
	}

	return matchingArtists
}

// CreateArtist create new artist
func (keeper Keeper) CreateArtist(ctx sdk.Context, name string, owner sdk.AccAddress) (types.Artist, sdk.Error) {
	artistID, err := keeper.GetArtistID(ctx)
	if err != nil {
		return types.Artist{}, err
	}

	artist := types.NewArtist(artistID, name, owner)

	keeper.SetArtist(ctx, artist)
	keeper.SetArtistID(ctx, artistID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateArtist,
			sdk.NewAttribute(types.AttributeKeyArtistID, fmt.Sprintf("%d", artistID)),
			sdk.NewAttribute(types.AttributeKeyArtistName, fmt.Sprintf("%s", name)),
			sdk.NewAttribute(types.AttributeKeyArtistOwner, fmt.Sprintf("%s", owner.String())),
		),
	)

	return artist, nil
}

// SetArtistImage create new artist
func (keeper Keeper) SetArtistImage(ctx sdk.Context, artistID uint64, image types.ArtistImage, owner sdk.AccAddress) sdk.Error {
	artist, ok := keeper.GetArtist(ctx, artistID)
	if !ok {
		return types.ErrUnknownArtist(keeper.codespace, "unknown artist")
	}

	if artist.Owner.String() != owner.String() {
		return types.ErrUnknownOwner(keeper.codespace, "unknown owner")
	}

	artist.Image = image

	keeper.SetArtist(ctx, artist)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetArtistImage,
			sdk.NewAttribute(types.AttributeKeyArtistID, fmt.Sprintf("%d", artistID)),
			sdk.NewAttribute(types.AttributeKeyArtistName, fmt.Sprintf("%s", artist.Name)),
			sdk.NewAttribute(types.AttributeKeyArtistImage, fmt.Sprintf("%s", artist.Image.CID)),
			sdk.NewAttribute(types.AttributeKeyArtistOwner, fmt.Sprintf("%s", artist.Owner.String())),
		),
	)

	return nil
}
