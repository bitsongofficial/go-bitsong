package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, acc := range data.Accounts {
		k.RegisterHandle(ctx, acc.GetAddress(), acc.Handle)
	}

	return []abci.ValidatorUpdate{}
}

/*func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	return GenesisState{
		Artists: k.GetAllArtists(ctx),
	}
}
*/
