package types

import "cosmossdk.io/collections"

const (
	ModuleName = "drop"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	MaxRulesPerDrop = 5
	MaxRuleIDLength = 20
)

var (
	DropsPrefix = collections.NewPrefix(0)
	RulesPrefix = collections.NewPrefix(1)
)
