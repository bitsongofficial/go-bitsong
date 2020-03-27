package track

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bitsongofficial/go-bitsong/x/track/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/track/client/rest"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the 20-transfer appmodulebasic
type AppModuleBasic struct{}

// Name implements AppModuleBasic interface
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec implements AppModuleBasic interface
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// transfer module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the track module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
	var data GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes implements AppModuleBasic interface
func (a AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd implements AppModuleBasic interface
func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd implements AppModuleBasic interface
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule represents the AppModule for this module
type AppModule struct {
	AppModuleBasic
	keeper     Keeper
	bankKeeper bank.Keeper
}

// NewAppModule creates a new 20-transfer module
func NewAppModule(keeper Keeper, bankKeeper bank.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		bankKeeper:     bankKeeper,
	}
}

// RegisterInvariants implements the AppModule interface
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
}

// Route implements the AppModule interface
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler implements the AppModule interface
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute implements the AppModule interface
func (AppModule) QuerierRoute() string {
	return ModuleName
}

// NewQuerierHandler implements the AppModule interface
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the track module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.bankKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the ibc
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the track module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the track module. It returns no validator
// updates.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
