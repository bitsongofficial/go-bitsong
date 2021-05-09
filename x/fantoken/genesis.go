package token

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/ledger/x/fantoken/keeper"
	"github.com/bitsongofficial/ledger/x/fantoken/types"
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

	// assert the denom exists
	if !k.HasDenom(ctx, data.Params.IssuePrice.Denom) {
		panic(fmt.Sprintf("Token %s does not exist", data.Params.IssuePrice.Denom))
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
		Tokens: []types.FanToken{types.GetNativeToken()},
	}
}
