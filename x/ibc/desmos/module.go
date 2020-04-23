package desmos

import (
	"encoding/json"

	"github.com/bitsongofficial/go-bitsong/x/ibc/desmos/client/cli"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the Desmos IBC AppModuleBasic
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
	return cdc.MustMarshalJSON(DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the ibc desmos module.
func (AppModuleBasic) ValidateGenesis(_ codec.JSONMarshaler, _ json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes implements AppModuleBasic interface
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	// TODO
}

// GetTxCmd implements AppModuleBasic interface
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetIBCDesmosTxCommand(cdc)
}

// GetQueryCmd implements AppModuleBasic interface
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// TODO
	return nil
}

// AppModule represents the AppModule for this module
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new 20-transfer module
func NewAppModule(k Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

// RegisterInvariants implements the AppModule interface
func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO
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
	return QuerierRoute
}

// NewQuerierHandler implements the AppModule interface
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

// InitGenesis performs genesis initialization for the ibc desmos module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {

}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// Implement IBCModule callbacks
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order channelexported.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	// TODO: Enforce ordering, currently relayers use ORDERED channels

	// Require portID is the portID the desmos module is bound to
	boundPort := am.keeper.GetPort(ctx)
	if boundPort != portID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if version != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	// TODO: escrow
	return nil
}

func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order channelexported.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version,
	counterpartyVersion string,
) error {
	// TODO: Enforce ordering, currently relayers use ORDERED channels

	if counterparty.PortID != PortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "counterparty has invalid portid. expected: %s, got %s", PortID, counterparty.PortID)
	}

	if version != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	}

	if counterpartyVersion != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	// TODO: escrow
	return nil
}

func (am AppModule) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyVersion string) error {
	if counterpartyVersion != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	}
	return nil
}

func (am AppModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (am AppModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (am AppModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (am AppModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
	return nil, nil
}

func (am AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledment []byte) (*sdk.Result, error) {
	return nil, nil
}

func (am AppModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
	return nil, nil
}
