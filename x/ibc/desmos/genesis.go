package desmos

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets ibc desmos information for genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, state GenesisState) {
	// desmos module binds to the desmos port on InitChain
	// and claims the returned capability
	err := keeper.BindPort(ctx, state.PortID)
	if err != nil {
		panic(fmt.Sprintf("could not claim port capability: %v", err))
	}
}

// ExportGenesis exports transfer module's portID into its geneis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	portID := keeper.GetPort(ctx)

	return GenesisState{
		PortID: portID,
	}
}
