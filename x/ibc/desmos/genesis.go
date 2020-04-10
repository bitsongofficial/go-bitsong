package desmos

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets ibc posts information for genesis
func InitGenesis(ctx sdk.Context, keeper Keeper) {
	// posts module binds to the posts port on InitChain
	// and claims the returned capability
	err := keeper.BindPort(ctx, PortID)
	if err != nil {
		panic(fmt.Sprintf("could not claim port capability: %v", err))
	}
}
