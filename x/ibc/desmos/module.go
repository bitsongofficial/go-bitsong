package desmos

import (
	"encoding/json"

	"github.com/bitsongofficial/go-bitsong/x/ibc/desmos/client/cli"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
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

}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// transfer module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return nil
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
