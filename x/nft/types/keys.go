package types

import "cosmossdk.io/collections"

const (
	ModuleName = "nft"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	CollectionsPrefix      = collections.NewPrefix(0)
	SupplyPrefix           = collections.NewPrefix(1)
	NFTsPrefix             = collections.NewPrefix(2)
	NFTsByCollectionPrefix = collections.NewPrefix(3)
	NFTsByOwnerPrefix      = collections.NewPrefix(4)
)
