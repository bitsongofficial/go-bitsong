package types

import "cosmossdk.io/collections"

const (
	ModuleName = "nft"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	CollectionsPrefix = collections.NewPrefix(0)
	SupplyPrefix      = collections.NewPrefix(1)
)
