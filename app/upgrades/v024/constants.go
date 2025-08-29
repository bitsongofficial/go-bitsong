package v024

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	icq "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
)

const (
	UpgradeName = "v024"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV023UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{Added: []string{}, Deleted: []string{icq.ModuleName}},
}
