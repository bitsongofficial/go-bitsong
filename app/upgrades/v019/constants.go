package v019

import (
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v019"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV019UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
