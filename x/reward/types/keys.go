package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "reward"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

// Keys for reward store
// Items are stored with the following key: values
//
// - 0x10<accAddr_Bytes>: Reward
//
var (
	RewardsKeyPrefix = []byte{0x10}
)

func RewardKey(accAddr sdk.AccAddress) []byte {
	return append(RewardsKeyPrefix, accAddr.Bytes()...)
}
