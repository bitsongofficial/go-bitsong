package track

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

// GenesisState - all track state that must be provided at genesis
type GenesisState struct {
	StartingTrackID uint64 `json:"starting_track_id"`
	Tracks          Tracks `json:"tracks"`
}

// NewGenesisState creates a new genesis state for the track module
func NewGenesisState(startingTrackID uint64) GenesisState {
	return GenesisState{
		StartingTrackID: startingTrackID,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingTrackID: 1,
	}
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// ValidateGenesis validates the given genesis state and returns an error if something is invalid
func ValidateGenesis(data GenesisState) error {
	// TODO: add validation
	/*for _, record := range data.Artists {
		if err := record.Validate(); err != nil {
			return err
		}
	}*/

	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetTrackID(ctx, data.StartingTrackID)

	for _, track := range data.Tracks {
		k.SetTrack(ctx, track)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingTrackID, _ := k.GetTrackID(ctx)
	// TODO: export only verified tracks?
	tracks := k.GetTracksFiltered(ctx, sdk.AccAddress{}, types.StatusVerified, 0)

	return GenesisState{
		StartingTrackID: startingTrackID,
		Tracks:          tracks,
	}
}
