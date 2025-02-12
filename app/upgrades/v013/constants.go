package v013

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
)

const (
	UpgradeName = "v013"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV13UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
