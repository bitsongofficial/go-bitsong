package types

import "fmt"

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	params Params,
	metadata []Metadata, lastMetadataId uint64,
	nfts []NFT, lastNftId uint64,
	collections []Collection, lastCollectionId uint64,
) GenesisState {
	return GenesisState{
		Params:           params,
		Metadata:         metadata,
		LastMetadataId:   lastMetadataId,
		Nfts:             nfts,
		Collections:      collections,
		LastCollectionId: lastCollectionId,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	for _, meta := range data.Metadata {
		if meta.Id > data.LastMetadataId {
			return fmt.Errorf("invalid metadata id: %d", meta.Id)
		}
	}

	for _, nft := range data.Nfts {
		if nft.MetadataId > data.LastMetadataId {
			return fmt.Errorf("invalid metadata id: %d", nft.MetadataId)
		}
	}

	for _, collection := range data.Collections {
		if collection.Id > data.LastCollectionId {
			return fmt.Errorf("invalid metadata id: %d", collection.Id)
		}
		if collection.Name == "" {
			return fmt.Errorf("invalid collection name")
		}
	}

	return nil
}
