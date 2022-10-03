package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores the genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	// initialize params
	k.SetParamSet(ctx, data.Params)

	// initialize metadata
	for _, metadata := range data.Metadata {
		k.SetMetadata(ctx, metadata)
	}

	// initialize last metadata ids
	for _, info := range data.LastMetadataIds {
		k.SetLastMetadataId(ctx, info.CollId, info.LastMetadataId)
	}

	// initialize nfts
	for _, nft := range data.Nfts {
		k.SetNFT(ctx, nft)
	}

	// initialize collections
	for _, collection := range data.Collections {
		k.SetCollection(ctx, collection)
	}
	k.SetLastCollectionId(ctx, data.LastCollectionId)
}

// ExportGenesis outputs the genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:           k.GetParamSet(ctx),
		Metadata:         k.GetAllMetadata(ctx),
		LastMetadataIds:  k.GetAllLastMetadataIds(ctx),
		Nfts:             k.GetAllNFTs(ctx),
		Collections:      k.GetAllCollections(ctx),
		LastCollectionId: k.GetLastCollectionId(ctx),
	}
}
