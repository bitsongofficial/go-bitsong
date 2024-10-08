package v017

import (
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	merkledroptypes "github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v017"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV17UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Deleted: []string{merkledroptypes.ModuleName},
	},
}
