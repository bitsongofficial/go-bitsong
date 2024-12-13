package v018

import (
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	merkledroptypes "github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
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
		Deleted: []string{merkledroptypes.ModuleName},
	},
}
