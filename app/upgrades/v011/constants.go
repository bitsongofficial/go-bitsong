package v011

import (
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v011"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV11UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{fantokentypes.ModuleName, "merkledrop"},
	},
}
