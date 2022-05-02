package nft

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/keeper"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params:            types.DefaultParams(),
		Metadata:          []types.Metadata{},
		LastMetadataId:    0,
		Nfts:              []types.NFT{},
		LastNftId:         0,
		Collections:       []types.Collection{},
		LastCollectionId:  0,
		CollectionRecords: []types.CollectionRecord{},
	}
}

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	// initialize params
	k.SetParamSet(ctx, data.Params)

	// initialize metadata
	for _, metadata := range data.Metadata {
		k.SetMetadata(ctx, metadata)
	}
	k.SetLastMetadataId(ctx, data.LastMetadataId)

	// initialize nfts
	for _, nft := range data.Nfts {
		k.SetNFT(ctx, nft)
	}
	k.SetLastNftId(ctx, data.LastNftId)

	// initialize collections
	for _, collection := range data.Collections {
		k.SetCollection(ctx, collection)
	}
	k.SetLastCollectionId(ctx, data.LastCollectionId)

	// initialize collection records
	for _, record := range data.CollectionRecords {
		k.SetCollectionNftRecord(ctx, record.CollectionId, record.NftId)
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:            k.GetParamSet(ctx),
		Metadata:          k.GetAllMetadata(ctx),
		LastMetadataId:    k.GetLastMetadataId(ctx),
		Nfts:              k.GetAllNFTs(ctx),
		LastNftId:         k.GetLastNftId(ctx),
		Collections:       k.GetAllCollections(ctx),
		LastCollectionId:  k.GetLastCollectionId(ctx),
		CollectionRecords: k.GetAllCollectionNftRecords(ctx),
	}
}
