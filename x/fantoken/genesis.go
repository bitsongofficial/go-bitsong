package fantoken

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	// init tokens
	for _, token := range data.Tokens {
		if err := k.AddFanToken(ctx, token); err != nil {
			panic(err.Error())
		}
	}

	for _, coin := range data.BurnedCoins {
		k.AddBurnCoin(ctx, coin)
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var tokens []types.FanToken
	for _, token := range k.GetFanTokens(ctx, nil) {
		t := token.(*types.FanToken)
		tokens = append(tokens, *t)
	}
	return &types.GenesisState{
		Params:      k.GetParamSet(ctx),
		Tokens:      tokens,
		BurnedCoins: k.GetAllBurnCoin(ctx),
	}
}

// DefaultGenesisState returns the default genesis state for testing
func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params: types.DefaultParams(),
		Tokens: []types.FanToken{},
	}
}
