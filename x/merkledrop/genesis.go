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
		Indexes:          []*types.Indexes{},
	}
}

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	// initialize merkledrops
	for _, md := range data.Merkledrops {
		k.SetMerkleDrop(ctx, md)
	}

	// set last merkledrop id
	k.SetLastMerkleDropId(ctx, data.LastMerkledropId)

	// TODO: set index claims
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		LastMerkledropId: k.GetLastMerkleDropId(ctx),
		Merkledrops:      k.GetAllMerkleDrops(ctx),
		Indexes:          k.GetAllIndexes(ctx),
	}
}
