package track

import (
	"fmt"
	"github.com/BitSongOfficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - track genesis state
type GenesisState struct {
	Pool            Pool         `json:"pool"`
	Params          types.Params `json:"params"`
	StartingTrackID uint64       `json:"starting_id"`
	Tracks          Tracks       `json:"tracks"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(pool Pool, startingTrackID uint64, params types.Params) GenesisState {
	return GenesisState{
		Pool:            pool,
		StartingTrackID: startingTrackID,
		Params:          params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Pool:            types.InitialPool(),
		StartingTrackID: DefaultStartingSongID,
		Params:          types.DefaultParams(),
	}
}

// ValidateGenesis validates genesis state
func ValidateGenesis(data GenesisState) error {
	if data.Params.PlayTax.IsNegative() || data.Params.PlayTax.GT(sdk.OneDec()) {
		return fmt.Errorf("track parameter PlayTax should non-negative and "+
			"less than one, is %s", data.Params.PlayTax.String())
	}

	return data.Pool.ValidateGenesis()
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetFeePlayPool(ctx, data.Pool)

	err := keeper.SetInitialTrackID(ctx, data.StartingTrackID)
	if err != nil {
		panic(err)
	}
	/*for _, track := range data.Tracks {
		keeper.AddSong(ctx, *track)
	}*/

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	startingSongID, _ := keeper.PeekCurrentTrackID(ctx)
	return NewGenesisState(Pool{}, startingSongID, keeper.GetParams(ctx))
}
