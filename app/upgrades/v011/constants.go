package v011

import (
	"github.com/bitsongofficial/go-bitsong/v018/app/upgrades"
	fantokentypes "github.com/bitsongofficial/go-bitsong/v018/x/fantoken/types"
	merkledroptypes "github.com/bitsongofficial/go-bitsong/v018/x/merkledrop/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	UpgradeName = "v011"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV11UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{fantokentypes.ModuleName, merkledroptypes.ModuleName},
	},
}
