package candymachine

import (
	"github.com/bitsongofficial/go-bitsong/x/candymachine/keeper"
	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params: types.DefaultParams(),
	}
}

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	// initialize params
	k.SetParamSet(ctx, data.Params)

	for _, machine := range data.Candymachines {
		k.SetCandyMachine(ctx, machine)
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:        k.GetParamSet(ctx),
		Candymachines: k.GetAllCandyMachines(ctx),
	}
}
