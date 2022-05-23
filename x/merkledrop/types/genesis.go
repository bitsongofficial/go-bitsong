package types

import (
	"fmt"
)

func NewGenesisState(lastMdId uint64, mds []Merkledrop, indexes []*Indexes) GenesisState {
	return GenesisState{
		LastMerkledropId: lastMdId,
		Merkledrops:      mds,
		Indexes:          indexes,
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, md := range data.Merkledrops {
		if md.Id > data.LastMerkledropId {
			return fmt.Errorf("invalid merlkedrop id: %d", md.Id)
		}
	}

	for _, i := range data.Indexes {
		if i.MerkledropId > data.LastMerkledropId {
			return fmt.Errorf("invalid index merkledrop_id: %d", i.MerkledropId)
		}
	}

	return nil
}
