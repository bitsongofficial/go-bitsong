package release

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, release := range data.Releases {
		k.SetRelease(ctx, release)
		k.SetReleaseForCreator(ctx, release)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	return GenesisState{
		Releases: k.GetAllReleases(ctx),
	}
}
