package v015

import (
	"github.com/bitsongofficial/go-bitsong/v018/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v015"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV15UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
