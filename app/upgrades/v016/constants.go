package v016

import (
	"github.com/bitsongofficial/go-bitsong/v018/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v016"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV16UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
