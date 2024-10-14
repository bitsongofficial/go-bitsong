package v013

import (
	"github.com/bitsongofficial/go-bitsong/v018/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v013"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV13UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
