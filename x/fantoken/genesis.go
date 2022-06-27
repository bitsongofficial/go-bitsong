package fantoken

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := data.Validate(); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	// init fan tokens
	for _, fantoken := range data.FanTokens {
		if err := k.AddFanToken(ctx, &fantoken); err != nil {
			panic(err.Error())
		}
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:    k.GetParamSet(ctx),
		FanTokens: k.GetFanTokens(ctx, nil),
	}
}
