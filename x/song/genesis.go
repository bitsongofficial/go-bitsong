package song

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - song genesis state
type GenesisState struct {
	StartingSongId uint64 `json:"starting_id"`
	Songs          Songs  `json:"songs"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(startingSongId uint64) GenesisState {
	return GenesisState{
		StartingSongId: startingSongId,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		StartingSongId: DefaultStartingSongId,
	}
}

// Checks whether 2 GenesisState structs are equivalent.
/*func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := MsgCdc.MustMarshalBinaryBare(data)
	b2 := MsgCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}*/

// Returns if a GenesisState is empty or has data in it
/*func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}*/

// ValidateGenesis validates genesis state
func ValidateGenesis(data GenesisState) error {
	return nil
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	err := keeper.SetInitialSongId(ctx, data.StartingSongId)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}
	/*for _, song := range data.Songs {
		keeper.AddSong(ctx, *song)
	}*/

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	startingSongId, _ := keeper.PeekCurrentSongId(ctx)
	//orders := keeper.GetOrdersFiltered(ctx, nil, "", "", 0)
	return GenesisState{
		StartingSongId: startingSongId,
		//Orders:          orders,
	}
}
/*func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState()
}*/