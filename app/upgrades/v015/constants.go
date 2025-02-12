package v015

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
)

const (
	UpgradeName = "v015"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV15UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
