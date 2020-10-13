package auth

import (
	"encoding/json"
	"github.com/bitsongofficial/go-bitsong/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the content module.
type AppModuleBasic struct{}

// Name returns the content module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the content module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	*CosmosModuleCdc = *ModuleCdc
}

// DefaultGenesis returns default genesis state as raw bytes for the content
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	//return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
	return CosmosAppModuleBasic{}.DefaultGenesis()
}

// ValidateGenesis performs genesis state validation for the content module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	/*var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)*/
	return CosmosAppModuleBasic{}.ValidateGenesis(bz)
}

// RegisterRESTRoutes registers the REST routes for the content module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	//rest.RegisterRoutes(ctx, rtr)
	CosmosAppModuleBasic{}.RegisterRESTRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the content module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the content module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the content module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule CosmosAppModule
	authKeeper      AccountKeeper

	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(ak AccountKeeper, bacck Keeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: NewCosmosAppModule(ak),
		keeper:          bacck,
	}
}

// Name returns the content module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the content module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the content module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the content module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the content module's querier route name.
func (am AppModule) QuerierRoute() string {
	return am.cosmosAppModule.QuerierRoute()
}

// NewQuerierHandler returns the content module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return am.cosmosAppModule.NewQuerierHandler()
}

// InitGenesis performs genesis initialization for the content module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the content
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx)
}

// BeginBlock returns the begin blocker for the content module.
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, rbb)
}

// EndBlock returns the end blocker for the content module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.cosmosAppModule.EndBlock(ctx, rbb)
}
