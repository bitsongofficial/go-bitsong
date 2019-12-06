package distribution

import (
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
)

const (
	ModuleName = distr.ModuleName
	StoreKey                         = distr.StoreKey
	DefaultParamspace                = distr.DefaultParamspace
	DefaultCodespace                 = distr.DefaultCodespace
	RouterKey                        = distr.RouterKey
)

var (
	NewKeeper                                  = distr.NewKeeper
	RegisterCodec = distr.RegisterCodec
	CosmosModuleCdc = distr.ModuleCdc
	NewCosmosAppModule    = distr.NewAppModule
	ProposalHandler = distr.ProposalHandler
	NewCommunityPoolSpendProposalHandler       = distr.NewCommunityPoolSpendProposalHandler
)

type (
	CosmosAppModule = distr.AppModule
	CosmosAppModuleBasic = distr.AppModuleBasic
	Keeper                                 = distr.Keeper
)