package v020

import (
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v18"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV182UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
