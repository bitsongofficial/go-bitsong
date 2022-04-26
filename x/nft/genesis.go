package nft

import "github.com/bitsongofficial/go-bitsong/x/nft/types"

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
