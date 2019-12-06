package mint

import (
	mintCosmos "github.com/cosmos/cosmos-sdk/x/mint"
)

const (
	ModuleName                       = mintCosmos.ModuleName
	StoreKey                         = mintCosmos.StoreKey
	DefaultParamspace                = mintCosmos.DefaultParamspace
)

var(
	NewCosmosAppModule    = mintCosmos.NewAppModule
	NewKeeper                                  = mintCosmos.NewKeeper

)

type (
	CosmosAppModule = mintCosmos.AppModule
	CosmosAppModuleBasic = mintCosmos.AppModuleBasic
	Keeper = mintCosmos.Keeper
)