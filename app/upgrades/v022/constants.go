package v022

import (
	store "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"
)

const (
	UpgradeName = "v022"
	PatchVal1   = "bitsongvaloper1qxw4fjged2xve8ez7nu779tm8ejw92rv0vcuqr"
	PatchVal2   = "bitsongvaloper1xnc32z84cc9vwftvv4w0v02a2slug3tjt6qyct"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV022UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{Added: []string{
		protocolpooltypes.StoreKey,
	}, Deleted: []string{}},
}

type ConditionalJSON struct {
	PatchDelegationCount     uint
	PatchedHistRewards       []distrtypes.ValidatorHistoricalRewardsRecord
	ZeroSharesDelegation     []ZeroSharesDelegation
	PatchedDelegation        []PatchedDelegation
	NilDelegationCalculation []NilDelegationCalculation
	DistSlashStore           DistrSlashObject
}

type DistrSlashObject struct {
	SlashEventCount uint64               `json:"total_slashes"`
	DistrSlashEvent []map[string][]Slash `json:"events"`
}
type DistrSlashEvent struct {
	Val             string  `json:"val_addr"`
	SlashEventCount uint64  `json:"total"`
	Slashes         []Slash `json:"slash_events"`
}
type Slash struct {
	Height   uint64 `json:"height"`
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
type NilDelegationCalculation struct {
	OperatorAddress  string `json:"val_addr"`
	DelegatorAddress string `json:"del_addr"`
}
