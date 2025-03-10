package v022

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

const (
	UpgradeName = "v022"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV022UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{Added: []string{}, Deleted: []string{}},
}

type ConditionalJSON struct {
	PatchDelegationCount uint
	PatchedHistRewards   []distrtypes.ValidatorHistoricalRewardsRecord
	ZeroSharesDelegation []ZeroSharesDelegation
	PatchedDelegation    []PatchedDelegation
	DistSlashStore       DistrSlashObject
}

type DistrSlashObject struct {
	SlashEventCount uint64            `json:"total_slashes"`
	DistrSlashEvent []DistrSlashEvent `json:"events"`
}
type DistrSlashEvent struct {
	Val             string  `json:"val_addr"`
	SlashEventCount uint64  `json:"total"`
	Slashes         []Slash `json:"slash_events"`
}
type Slash struct {
	Fraction string `json:"fraction"`
	Period   uint64 `json:"period"`
}

type ZeroSharesDelegation struct {
	OperatorAddress  string `json:"val_addr"`
	DelegatorAddress string `json:"del_addr"`
}
type PatchedDelegation struct {
	OperatorAddress   string `json:"val_addr"`
	DelegatorAddress  string `json:"del_addr"`
	PatchedDelegation string `json:"patch"`
}
