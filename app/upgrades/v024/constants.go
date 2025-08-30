package v024

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
)

const (
	UpgradeName = "v024"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV024Upgrade,
	StoreUpgrades:        store.StoreUpgrades{Added: []string{}, Deleted: []string{"interchainquery"}},
}
