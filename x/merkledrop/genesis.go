package merkledrop

import (
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/keeper"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		LastMerkledropId: 0,
		Merkledrops:      []types.Merkledrop{},
	}
}

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	// initialize merkledrops
	if data.LastMerkledropId == 0 {
		k.CreateModuleAccount(ctx)
	}

	for _, md := range data.Merkledrops {
		k.SetMerkleDrop(ctx, md)
	}

	// set last merkledrop id
	k.SetLastMerkleDropId(ctx, data.LastMerkledropId)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	// TODO: export index claimed

	return &types.GenesisState{
		LastMerkledropId:     k.GetLastMerkleDropId(ctx),
		Merkledrops:          k.GetAllMerkleDrops(ctx),
		ModuleAccountBalance: k.GetModuleAccountBalance(ctx),
	}
}