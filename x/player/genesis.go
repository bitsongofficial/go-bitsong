package player

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, p := range data.Players {
		k.SetPlayer(ctx, p)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	// TODO: add players iteration and export genesis
	return GenesisState{
		Players: []Player{},
	}
}
