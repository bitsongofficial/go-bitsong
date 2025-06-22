package v023

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
)

const (
	UpgradeName = "v023"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV023UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{Added: []string{}, Deleted: []string{}},
}
