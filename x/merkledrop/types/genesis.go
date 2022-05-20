package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGenesisState(lastMdId uint64, mds []Merkledrop) GenesisState {
	return GenesisState{
		LastMerkledropId:     lastMdId,
		Merkledrops:          mds,
		ModuleAccountBalance: sdk.Coins{},
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, md := range data.Merkledrops {
		if md.Id > data.LastMerkledropId {
			return fmt.Errorf("invalid merlkedrop id: %d", md.Id)
		}
	}

	return nil
}
