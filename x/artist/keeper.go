package artist

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

// Governance Keeper
type Keeper struct {
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

// NewKeeper returns an artist keeper.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType,
) Keeper {
	// TODO:
	// need router.seal() ???

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

// Set the artist ID
func (keeper Keeper) setArtistID(ctx sdk.Context, artistID uint64) {
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
// status will filter proposals by status
// numLatest will fetch a specified number of the most recent proposals, or 0 for all proposals
func (keeper Keeper) GetArtistsFiltered(ctx sdk.Context, ownerAddr sdk.AccAddress, status types.ArtistStatus, numLatest uint64) []types.Artist {

	maxArtistID, err := keeper.GetArtistID(ctx)
	if err != nil {
		return []types.Artist{}
	}

	matchingArtists := []types.Artist{}

	if numLatest == 0 {
		numLatest = maxArtistID
	}

	for artistID := maxArtistID - numLatest; artistID < maxArtistID; artistID++ {
		artist, ok := keeper.GetArtist(ctx, artistID)
		if !ok {
			continue
		}

		if types.ValidArtistStatus(status) && artist.Status != status {
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
func (keeper Keeper) CreateArtist(ctx sdk.Context, meta types.Meta, owner sdk.AccAddress) (types.Artist, sdk.Error) {
	artistID, err := keeper.GetArtistID(ctx)
	if err != nil {
		return types.Artist{}, err
	}

	artist := types.NewArtist(artistID, meta, owner)

	keeper.SetArtist(ctx, artist)
	keeper.setArtistID(ctx, artistID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateArtist,
			sdk.NewAttribute(types.AttributeKeyArtistID, fmt.Sprintf("%d", artistID)),
			sdk.NewAttribute(types.AttributeKeyArtistMeta, fmt.Sprintf("%s", meta.String())),
			sdk.NewAttribute(types.AttributeKeyArtistOwner, fmt.Sprintf("%s", owner.String())),
		),
	)

	return artist, nil
}
