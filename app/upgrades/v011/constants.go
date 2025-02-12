package v011

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

const (
	UpgradeName = "v011"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV11UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{fantokentypes.ModuleName, "merkledrop"},
	},
}
