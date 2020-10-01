package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "release"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

// Keys for release store
// Items are stored with the following key: values
//
// - 0x00<releaseID_Bytes>: Release
// - 0x10<accAddr_Bytes><releaseID_Bytes>: releaseID
var (
	CounterKey                 = []byte{0x00}
	ReleaseKeyPrefix           = []byte{0x10}
	ReleaseForCreatorKeyPrefix = []byte{0x20}
)

func GetReleaseKey(releaseID string) []byte {
	return append(ReleaseKeyPrefix, []byte(releaseID)...)
}

func GetReleaseForCreatorKey(address sdk.AccAddress) []byte {
	return append(ReleaseForCreatorKeyPrefix, address.Bytes()...)
}

func ReleaseAddressKey(address sdk.AccAddress, releaseID string) []byte {
	return append(GetReleaseForCreatorKey(address), []byte(releaseID)...)
}
