package distribution

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
)

type OverrideDistrModule struct {
	distr.AppModule
	k OverrideDistrKeeper
}

func NewOverrideDistrModule(appModule distr.AppModule, keeper OverrideDistrKeeper) OverrideDistrModule {
	return OverrideDistrModule{
		AppModule: appModule,
		k:         keeper,
	}
}

// module name
func (am OverrideDistrModule) Name() string {
	return am.AppModule.Name()
}

// register invariants
func (am OverrideDistrModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.AppModule.RegisterInvariants(ir)
}

// module querier route name
func (am OverrideDistrModule) Route() string {
	return am.AppModule.Route()
}

// module handler
func (am OverrideDistrModule) NewHandler() sdk.Handler {
	return am.AppModule.NewHandler()
}

// module querier route name
func (am OverrideDistrModule) QuerierRoute() string { return am.AppModule.QuerierRoute() }

// module querier
func (am OverrideDistrModule) NewQuerierHandler() sdk.Querier { return am.AppModule.NewQuerierHandler() }

// module init-genesis
func (am OverrideDistrModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.AppModule.InitGenesis(ctx, data)
}

// module export genesis
func (am OverrideDistrModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.AppModule.ExportGenesis(ctx)
}

// module begin-block
func (am OverrideDistrModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	BeginBlocker(ctx, rbb, am.k)
}

// module end-block
func (am OverrideDistrModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.AppModule.EndBlock(ctx, rbb)
}
