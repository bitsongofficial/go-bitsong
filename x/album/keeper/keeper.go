package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   sdk.StoreKey       // The (unexposed) keys used to access the stores from the Context.
	cdc        *codec.Codec       // The codec for binary encoding/decoding.
	codespace  sdk.CodespaceType  // Reserved codespace
	ak         auth.AccountKeeper // Cosmos-SDK Account Keeper
	Sk         supply.Keeper      // Cosmos-SDK Supply Keeper
	paramSpace params.Subspace
}

// NewKeeper returns an album keeper.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType, ak auth.AccountKeeper, sk supply.Keeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		codespace:  codespace,
		ak:         ak,
		Sk:         sk,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

/****************************************
 * Album
 ****************************************/

// Set the album ID
func (keeper Keeper) SetAlbumID(ctx sdk.Context, albumID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(albumID)
	store.Set(types.AlbumIDKey, bz)
}

// GetAlbumID gets the highest album ID
func (keeper Keeper) GetAlbumID(ctx sdk.Context) (albumID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.AlbumIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis(keeper.codespace, "initial album ID hasn't been set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &albumID)
	return albumID, nil
}

// SetAlbum set an album to store
func (keeper Keeper) SetAlbum(ctx sdk.Context, album types.Album) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(album)
	store.Set(types.AlbumKey(album.AlbumID), bz)
}

// GetAlbum get Album from store by AlbumID
func (keeper Keeper) GetAlbum(ctx sdk.Context, albumID uint64) (album types.Album, ok bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.AlbumKey(albumID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &album)
	return album, true
}

// GetAlbumsFiltered get Albums from store by AlbumID
// status will filter albums by status
// numLatest will fetch a specified number of the most recent albums, or 0 for all albums
func (keeper Keeper) GetAlbumsFiltered(ctx sdk.Context, ownerAddr sdk.AccAddress, status types.AlbumStatus, numLatest uint64) []types.Album {

	maxAlbumID, err := keeper.GetAlbumID(ctx)
	if err != nil {
		return []types.Album{}
	}

	var matchingAlbums []types.Album

	if numLatest == 0 {
		numLatest = maxAlbumID
	}

	for albumID := maxAlbumID - numLatest; albumID < maxAlbumID; albumID++ {
		album, ok := keeper.GetAlbum(ctx, albumID)
		if !ok {
			continue
		}

		if album.Status.Valid() && album.Status != status {
			continue
		}

		if ownerAddr != nil && len(ownerAddr) != 0 && album.Owner.String() != ownerAddr.String() {
			continue
		}

		matchingAlbums = append(matchingAlbums, album)
	}

	return matchingAlbums
}

// CreateAlbum create new album
func (keeper Keeper) CreateAlbum(ctx sdk.Context, title string, albumType types.AlbumType, metadataUri string, owner sdk.AccAddress) (types.Album, sdk.Error) {
	albumID, err := keeper.GetAlbumID(ctx)
	if err != nil {
		return types.Album{}, err
	}

	submitTime := ctx.BlockHeader().Time

	album := types.NewAlbum(albumID, title, albumType, metadataUri, owner, submitTime)

	keeper.SetAlbum(ctx, album)
	keeper.SetAlbumID(ctx, albumID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateAlbum,
			sdk.NewAttribute(types.AttributeKeyAlbumID, fmt.Sprintf("%d", albumID)),
			sdk.NewAttribute(types.AttributeKeyAlbumTitle, fmt.Sprintf("%s", title)),
			sdk.NewAttribute(types.AttributeKeyAlbumOwner, fmt.Sprintf("%s", owner.String())),
		),
	)

	return album, nil
}

// SetAlbumStatus set Status of the Album {Nil, Verified, Rejected, Failed}
func (keeper Keeper) SetAlbumStatus(ctx sdk.Context, albumID uint64, status types.AlbumStatus) sdk.Error {
	album, ok := keeper.GetAlbum(ctx, albumID)
	if !ok {
		return types.ErrUnknownAlbum(keeper.codespace, "unknown album")
	}

	album.Status = status

	keeper.SetAlbum(ctx, album)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetAlbumStatus,
			sdk.NewAttribute(types.AttributeKeyAlbumID, fmt.Sprintf("%d", albumID)),
			sdk.NewAttribute(types.AttributeKeyAlbumTitle, fmt.Sprintf("%s", album.Title)),
			sdk.NewAttribute(types.AttributeKeyAlbumStatus, fmt.Sprintf("%s", album.Status.String())),
			sdk.NewAttribute(types.AttributeKeyAlbumOwner, fmt.Sprintf("%s", album.Owner.String())),
		),
	)

	return nil
}
