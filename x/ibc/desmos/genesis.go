package desmos

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState is currently only used to ensure that the InitGenesis gets run
// by the module manager
type GenesisState struct {
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

func DefaultGenesis() GenesisState {
	return GenesisState{
		Version: Version,
	}
}

// InitGenesis sets ibc posts information for genesis
func InitGenesis(ctx sdk.Context, keeper Keeper) {
	// posts module binds to the posts port on InitChain
	// and claims the returned capability
	err := keeper.BindPort(ctx, PortID)
	if err != nil {
		panic(fmt.Sprintf("could not claim port capability: %v", err))
	}
}
