package v010

import (
	store "cosmossdk.io/store/types"
	"cosmossdk.io/x/feegrant"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	"github.com/cosmos/cosmos-sdk/x/authz"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
)

const (
	UpgradeName       = "v010"
	CassiniMultiSig   = "bitsong12r2d9hhnd2ez4kgk63ar8m40vhaje8yaa94h8w"
	CassiniMintAmount = 9_656_879_130_000
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV10UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{authz.ModuleName, feegrant.ModuleName, packetforwardtypes.ModuleName},
	},
}
