package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "distributor"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

// Keys for distributor store
// Items are stored with the following key: values
//
// - 0x00<addr_Bytes>: Distributor
//
var (
	DistributorsKeyPrefix = []byte{0x00}
)

func DistributorKey(accaddr sdk.AccAddress) []byte {
	return append(DistributorsKeyPrefix, accaddr.Bytes()...)
}
