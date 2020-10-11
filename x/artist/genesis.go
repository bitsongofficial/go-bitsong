package artist

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, artist := range data.Artists {
		k.SetArtist(ctx, artist)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	return GenesisState{
		Artists: k.GetAllArtists(ctx),
	}
}
