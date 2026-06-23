package v024

import (
	store "cosmossdk.io/store/types"

	hyperlanetypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"

	"github.com/bitsongofficial/go-bitsong/app/upgrades"
)

const (
	UpgradeName = "v024"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV024UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			hyperlanetypes.ModuleName,
			warptypes.ModuleName,
		},
		Deleted: []string{},
	},
}
