package app

import (
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CosmosApp implements the common methods for a Cosmos SDK-based application
// specific blockchain.
type CosmosApp interface {
	// Name The assigned name of the app.
	Name() string

	// LegacyAmino The application types codec.
	// NOTE: This shoult be sealed before being returned.
	LegacyAmino() *codec.LegacyAmino

	// BeginBlocker Application updates every begin block.
	BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error)

	// EndBlocker Application updates every end block.
	EndBlocker(ctx sdk.Context) (sdk.EndBlock, error)

	// InitChainer Application update at chain (i.e app) initialization.
	InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error)

	// LoadHeight Loads the app at a given height.
	LoadHeight(height int64) error

	// ExportAppStateAndValidators Exports the state of the application for a genesis file.
	ExportAppStateAndValidators(
		forZeroHeight bool, jailAllowedAddrs []string,
	) (types.ExportedApp, error)

	// ModuleAccountAddrs All the registered module account addreses.
	ModuleAccountAddrs() map[string]bool
}
