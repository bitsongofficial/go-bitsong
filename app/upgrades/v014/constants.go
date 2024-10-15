package v014

import (
	"github.com/bitsongofficial/go-bitsong/v018/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v014"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV14UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
