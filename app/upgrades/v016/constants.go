package v016

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
)

const (
	UpgradeName = "v016"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV16UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
