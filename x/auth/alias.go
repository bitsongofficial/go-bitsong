package auth

import (
	"github.com/bitsongofficial/go-bitsong/x/auth/keeper"
	"github.com/bitsongofficial/go-bitsong/x/auth/types"
)

const (
	RouterKey = types.RouterKey
)

var (
	RegisterCodec = types.RegisterCodec
	NewHandler    = keeper.NewHandler
	NewKeeper     = keeper.NewKeeper
	//NewQuerier          = keeper.NewQuerier
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
