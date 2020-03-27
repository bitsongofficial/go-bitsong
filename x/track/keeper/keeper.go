package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey      sdk.StoreKey       // The (unexposed) keys used to access the stores from the Context.
	cdc           *codec.Codec       // The codec for binary encoding/decoding.
	codespace     string             // Reserved codespace
	stakingKeeper staking.Keeper     // Staking Keeper
	ak            auth.AccountKeeper // Cosmos-SDK Account Keeper
	Sk            supply.Keeper      // Cosmos-SDK Supply Keeper
	paramSpace    params.Subspace
}

// NewKeeper returns an track keeper.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace string, stakingKeeper staking.Keeper, ak auth.AccountKeeper, sk supply.Keeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		codespace:     codespace,
		stakingKeeper: stakingKeeper,
		ak:            ak,
		Sk:            sk,
		paramSpace:    paramSpace.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

/****************************************
 * Track
 ****************************************/

// Set the track ID
func (keeper Keeper) SetTrackID(ctx sdk.Context, trackID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(trackID)
	store.Set(types.TrackIDKey, bz)
}

// GetTrackID gets the highest track ID
func (keeper Keeper) GetTrackID(ctx sdk.Context) (trackID uint64, err error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.TrackIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis(keeper.codespace, "initial track ID hasn't been set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &trackID)
	return trackID, nil
}

// SetTrack set an track to store
func (keeper Keeper) SetTrack(ctx sdk.Context, track types.Track) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(track)
	store.Set(types.TrackKey(track.TrackID), bz)
}

// GetTrack get Track from store by TrackID
func (keeper Keeper) GetTrack(ctx sdk.Context, trackID uint64) (track types.Track, ok bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.TrackKey(trackID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &track)
	return track, true
}

// GetTracksFiltered get Tracks from store by TrackID
// status will filter tracks by status
// numLatest will fetch a specified number of the most recent tracks, or 0 for all tracks
func (keeper Keeper) GetTracksFiltered(ctx sdk.Context, ownerAddr sdk.AccAddress, status types.TrackStatus, numLatest uint64) []types.Track {

	maxTrackID, err := keeper.GetTrackID(ctx)
	if err != nil {
		return []types.Track{}
	}

	var matchingTracks []types.Track

	if numLatest == 0 {
		numLatest = maxTrackID
	}

	for trackID := maxTrackID - numLatest; trackID < maxTrackID; trackID++ {
		track, ok := keeper.GetTrack(ctx, trackID)
		if !ok {
			continue
		}

		if track.Status.Valid() && track.Status != status {
			continue
		}

		if ownerAddr != nil && len(ownerAddr) != 0 && track.Owner.String() != ownerAddr.String() {
			continue
		}

		matchingTracks = append(matchingTracks, track)
	}

	return matchingTracks
}

// CreateTrack create new track
func (keeper Keeper) CreateTrack(ctx sdk.Context,
	title, audio, image, duration string, hidden, explicit bool,
	genre, mood, artists, featuring, producers, description, copyright string,
	owner sdk.AccAddress,
) (types.Track, error) {
	trackID, err := keeper.GetTrackID(ctx)
	if err != nil {
		return types.Track{}, err
	}

	submitTime := ctx.BlockHeader().Time

	track := types.NewTrack(
		trackID,
		title,
		audio,
		image,
		duration,
		hidden,
		explicit,
		genre,
		mood,
		artists,
		featuring,
		producers,
		description,
		copyright,
		owner,
		submitTime,
	)

	keeper.SetTrack(ctx, track)
	keeper.SetTrackID(ctx, trackID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTrack,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", trackID)),
			sdk.NewAttribute(types.AttributeKeyTrackTitle, fmt.Sprintf("%s", title)),
			sdk.NewAttribute(types.AttributeKeyTrackOwner, fmt.Sprintf("%s", owner.String())),
		),
	)

	return track, nil
}

// SetTrackStatus set Status of the Track {Nil, Verified, Rejected, Failed}
func (keeper Keeper) SetTrackStatus(ctx sdk.Context, trackID uint64, status types.TrackStatus) error {
	track, ok := keeper.GetTrack(ctx, trackID)
	if !ok {
		return types.ErrUnknownTrack(keeper.codespace, "unknown track")
	}

	track.Status = status

	keeper.SetTrack(ctx, track)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetTrackStatus,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", trackID)),
		),
	)

	return nil
}
