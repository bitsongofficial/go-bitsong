package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
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

// NewKeeper returns an artist keeper.
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
 * Artist
 ****************************************/

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

// Set the artist ID
func (keeper Keeper) SetArtistID(ctx sdk.Context, artistID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(artistID)
	store.Set(types.ArtistIDKey, bz)
}

// SetArtist set an artist to store
func (keeper Keeper) SetArtist(ctx sdk.Context, artist types.Artist) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(artist)
	store.Set(types.ArtistKey(artist.ArtistID), bz)
}

// GetArtist get Artist from store by AlbumID
func (keeper Keeper) GetArtist(ctx sdk.Context, artistID uint64) (artist types.Artist, ok bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ArtistKey(artistID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &artist)
	return artist, true
}

// GetArtistsFiltered get Artists from store by AlbumID
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

func (keeper Keeper) PayFee(ctx sdk.Context, owner sdk.AccAddress, amt sdk.Coins) sdk.Error {
	// Get block time
	blockTime := ctx.BlockHeader().Time

	// Get account
	fromAcc := keeper.ak.GetAccount(ctx, owner)

	// Safe sub coins from account, verify the account has enough funds to pay for fees
	if _, hasNeg := fromAcc.GetCoins().SafeSub(amt); hasNeg {
		//return sdk.ErrInsufficientCoins(fmt.Sprintf("%s", fromAcc.GetCoins().String()))
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", fromAcc.GetCoins().String(), amt),
		)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := fromAcc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(amt); hasNeg {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", spendableCoins, amt),
		)
	}

	// Send fee from account to distribution module
	if err := keeper.Sk.SendCoinsFromAccountToModule(ctx, owner, distribution.ModuleName, amt); err != nil {
		return err
	}

	return nil
}

// CreateArtist create new artist
func (keeper Keeper) CreateArtist(ctx sdk.Context, name string, metadataUri string, owner sdk.AccAddress) (types.Artist, sdk.Error) {
	artistID, err := keeper.GetArtistID(ctx)
	if err != nil {
		return types.Artist{}, err
	}

	//////////////////////////////////////////
	// TODO: just for test, pay a fee to create a new artist
	//////////////////////////////////////////
	// feeAmt := sdk.Coins{sdk.NewCoin(btsg.BondDenom, sdk.NewInt(1000000))} // 1btsg = 1000000ubtsg
	// if err := keeper.PayFee(ctx, owner, feeAmt); err != nil {
	//	return types.Artist{}, err
	// }
	//////////////////////////////////////////

	submitTime := ctx.BlockHeader().Time

	artist := types.NewArtist(artistID, name, metadataUri, owner, submitTime)

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

// SetArtistStatus set Status of the Artist {Nil, Verified, Rejected, Failed}
func (keeper Keeper) SetArtistStatus(ctx sdk.Context, artistID uint64, status types.ArtistStatus) sdk.Error {
	artist, ok := keeper.GetArtist(ctx, artistID)
	if !ok {
		return types.ErrUnknownArtist(keeper.codespace, "unknown artist")
	}

	// TODO:
	// Interim moderator that can edit status on modules
	// This is only for testnet use and will be excluded then
	// if from.String() != util.ModeratorBech32AccAddr {
	//	return types.ErrUnknownModerator(keeper.codespace, "unknown moderator")
	// }

	artist.Status = status

	keeper.SetArtist(ctx, artist)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetArtistStatus,
			sdk.NewAttribute(types.AttributeKeyArtistID, fmt.Sprintf("%d", artistID)),
		),
	)

	return nil
}
