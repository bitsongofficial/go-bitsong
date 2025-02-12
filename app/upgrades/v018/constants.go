package v018

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

const (
	UpgradeName = "v18"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV18UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{crisistypes.ModuleName, consensustypes.ModuleName},
		Deleted: []string{"merkledrop"},
	},
}
