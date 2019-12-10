package distributor

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Distributors Distributors `json:"distributors"`
}

func NewGenesisState() GenesisState {
	return GenesisState{
		Distributors: Distributors{},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Distributors: Distributors{},
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
	for _, distr := range data.Distributors {
		k.SetDistributor(ctx, distr)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	/*tracks := k.GetTracksFiltered(ctx, sdk.AccAddress{}, types.StatusVerified, 0)

	return GenesisState{
		StartingTrackID: startingTrackID,
		Tracks:          tracks,
	}*/
	// TODO: add distributor export
	return GenesisState{
		Distributors: Distributors{},
	}
}
