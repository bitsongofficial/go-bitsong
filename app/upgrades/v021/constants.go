package v021

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
)

const (
	UpgradeName = "v021"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV021UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{icqtypes.StoreKey, ibcwasmtypes.StoreKey, ibchookstypes.StoreKey},
		Deleted: []string{},
	},
}
