package v021

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	cadancetypes "github.com/bitsongofficial/go-bitsong/x/cadance/types"
	smartaccounttypes "github.com/bitsongofficial/go-bitsong/x/smart-account/types"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
	wasmlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
)

const (
	UpgradeName = "v021"
	// MaximumUnauthenticatedGas for smart account transactions to verify the fee payer
	MaximumUnauthenticatedGas = uint64(120_000)
	// IsSmartAccountActive is used for the smart account circuit breaker, smartaccounts are deactivated for v25
	IsSmartAccountActive = false

	// CircuitBreakerController is a DAODAO address, used only to deactivate the smart account module
	// https://daodao.zone/dao/bitsong13hmdq0slwmff7sej79kfa8mgnx4rl46nj2fvmlgu6u32tz6vfqesdfq4vm/home
	CircuitBreakerController = "bitsong13hmdq0slwmff7sej79kfa8mgnx4rl46nj2fvmlgu6u32tz6vfqesdfq4vm"
)

var DefaultAllowedClients = []string{"07-tendermint", "09-localhost", wasmlctypes.Wasm}
var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV021UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			icqtypes.StoreKey,
			ibcwasmtypes.StoreKey,
			ibchookstypes.StoreKey,
			cadancetypes.StoreKey,
			smartaccounttypes.StoreKey,
		},
		Deleted: []string{},
	},
}
