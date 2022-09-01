package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores the genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	// initialize params
	k.SetParamSet(ctx, data.Params)

	for _, pad := range data.Launchpads {
		k.SetLaunchPad(ctx, pad)
	}

	for _, m := range data.MintableMetadataIds {
		k.SetMintableMetadataIds(ctx, m.CollectionId, m.MintableMetadataIds)
	}
}

// ExportGenesis outputs the genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:              k.GetParamSet(ctx),
		Launchpads:          k.GetAllLaunchPads(ctx),
		MintableMetadataIds: k.AllMintableMetadataIds(ctx),
	}
}
