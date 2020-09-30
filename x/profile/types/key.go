package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "profile"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

// Keys for account store
// Items are stored with the following key: values
//
// - 0x00<accAddr_Bytes>: Profile
var (
	ProfileKeyPrefix = []byte{0x00}
)

func GetProfileKey(accAddr sdk.AccAddress) []byte {
	return append(ProfileKeyPrefix, accAddr.Bytes()...)
}
