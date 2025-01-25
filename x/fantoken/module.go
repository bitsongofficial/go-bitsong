package fantoken

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/math"

	// simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModuleBasic = (*AppModule)(nil)
	// _ module.AppModuleSimulation = (*AppModule)(nil)
	_ module.HasGenesis = (*AppModule)(nil)

	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
	_ appmodule.AppModule       = AppModule{}
	_ module.AppModuleBasic     = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the fantoken module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the fantoken module's name.
func (AppModuleBasic) Name() string { return types.ModuleName }

// DefaultGenesis returns default genesis state as raw bytes for the fantoken module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the fantoken module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return data.Validate()
}

// RegisterRESTRoutes registers the REST routes for the fantoken module.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	// only grpc
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the fantoken module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the fantoken module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the fantoken module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterLegacyAminoCodec registers the fantoken module's types on the LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers interfaces and implementations of the fantoken module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// ____________________________________________________________________________

// AppModule implements an application module for the fantoken module.
type AppModule struct {
	AppModuleBasic

	keeper     keeper.Keeper
	bankKeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, bk types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		bankKeeper:     bk,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType is a marker function just indicates that this is a one-per-module type.
func (am AppModule) IsOnePerModuleType() {}

// Name returns the fantoken module's name.
func (AppModule) Name() string { return types.ModuleName }

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(&am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the fantoken module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns the message routing key for the fantoken module.
func (am AppModule) Route() string {
	return types.RouterKey
}

// QuerierRoute returns the fantoken module's querier route name.
func (AppModule) QuerierRoute() string { return types.RouterKey }

// InitGenesis performs genesis initialization for the fantoken module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, genesisState)
	return
}

// ExportGenesis returns the exported genesis state as raw bytes for the fantoken module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock performs a no-op.
func (AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// EndBlock returns the end blocker for the fantoken module. It returns no validator updates.
func (AppModule) EndBlock(_ context.Context) error {
	return nil
}

// GenerateGenesisState creates a randomized GenState of the valset module.
func (am AppModule) SimulatorGenesisState(simState *module.SimulationState) {
	ftDefaultGen := types.DefaultGenesisState()
	ftDefaultGen.Params.IssueFee = sdk.NewCoin("stake", math.NewInt(10000000))
	ftDefaultGenJson := simState.Cdc.MustMarshalJSON(ftDefaultGen)
	simState.GenState[types.ModuleName] = ftDefaultGenJson
}

// // GenerateGenesisState creates a randomized GenState of the module.
// func (a AppModule) GenerateGenesisState(_ *module.SimulationState) {}

// // RegisterStoreDecoder registers a decoder for module simulation testing
// func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {}
