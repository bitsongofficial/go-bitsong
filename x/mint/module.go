package mint

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cmint "github.com/cosmos/cosmos-sdk/x/mint"
)

var (
	_ module.AppModule = AppModule{}
)

// AppModule implements an application module for the mint module.
type AppModule struct {
	cmint.AppModule
	cmintKeeper cmint.Keeper
	keeper      Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(module cmint.AppModule, mintKeeper cmint.Keeper, keeper Keeper) AppModule {
	return AppModule{
		AppModule:   module,
		cmintKeeper: mintKeeper,
		keeper:      keeper,
	}
}

// module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.cmintKeeper, am.keeper)
}
