package artist

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/keeper"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
)

var (
	RegisterCodec       = types.RegisterCodec
	NewHandler          = keeper.NewHandler
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	ModuleCdc           = types.ModuleCdc
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
