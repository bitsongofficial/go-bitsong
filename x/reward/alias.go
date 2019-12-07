package reward

import (
	"github.com/bitsongofficial/go-bitsong/x/reward/keeper"
	"github.com/bitsongofficial/go-bitsong/x/reward/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	DefaultParamspace = keeper.DefaultParamspace
)

var (
	NewQuerier = keeper.NewQuerier

	// Codec
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
