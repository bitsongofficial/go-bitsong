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
// - 0x00<handle_Bytes>: Profile
// - 0x10<accAddr_Bytes>: Creator
var (
	ReleaseKeyPrefix = []byte{0x00}
	CreatorKeyPrefix = []byte{0x10}
)

func GetReleaseKey(releaseID string) []byte {
	return append(ReleaseKeyPrefix, []byte(releaseID)...)
}

func GerCreatorKey(address sdk.AccAddress) []byte {
	return append(CreatorKeyPrefix, address.Bytes()...)
}
