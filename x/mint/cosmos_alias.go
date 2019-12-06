package mint

import (
	mintCosmos "github.com/cosmos/cosmos-sdk/x/mint"
)

const (
	ModuleName                       = mintCosmos.ModuleName
	StoreKey                         = mintCosmos.StoreKey
	DefaultParamspace                = mintCosmos.DefaultParamspace

	EventTypeMint = mintCosmos.ModuleName
	AttributeKeyBondedRatio      = "bonded_ratio"
	AttributeKeyInflation        = "inflation"
	AttributeKeyAnnualProvisions = "annual_provisions"
)

var (
	NewKeeper             = mintCosmos.NewKeeper
	NewCosmosAppModule                         = mintCosmos.NewAppModule
)

type (
	CosmosAppModule = mintCosmos.AppModule
	CosmosAppModuleBasic = mintCosmos.AppModuleBasic
	Keeper = mintCosmos.Keeper
)