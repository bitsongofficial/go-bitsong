package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores the genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	// initialize params
	k.SetParamSet(ctx, data.Params)

	for _, machine := range data.Candymachines {
		k.SetCandyMachine(ctx, machine)
	}
}

// ExportGenesis outputs the genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:        k.GetParamSet(ctx),
		Candymachines: k.GetAllCandyMachines(ctx),
	}
}
