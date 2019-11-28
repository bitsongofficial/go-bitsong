package artist

import (
	"bytes"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all artist state that must be provided at genesis
type GenesisState struct {
	StartingArtistID uint64         `json:"starting_artist_id" yaml:"starting_artist_id"`
	Artists          []types.Artist `json:"artists" yaml:"artists"`
}

// NewGenesisState creates a new genesis state for the artist module
func NewGenesisState(startingArtistID uint64) GenesisState {
	return GenesisState{
		StartingArtistID: startingArtistID,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingArtistID: 1,
	}
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := types.ModuleCdc.MustMarshalBinaryBare(data)
	b2 := types.ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	return nil
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.setArtistID(ctx, data.StartingArtistID)

	for _, artist := range data.Artists {
		k.SetArtist(ctx, artist)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingArtistID, _ := k.GetArtistID(ctx)
	//artists := k.GetArtistsFiltered(ctx, types.StatusVerified, 0)
	artists := k.GetArtistsFiltered(ctx, types.StatusNil, 0)

	return GenesisState{
		StartingArtistID: startingArtistID,
		Artists:          artists,
	}
}
