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
// - 0x00<handle_Bytes>: Profile
// - 0x10<accAddr_Bytes>: Handle
var (
	ProfileKeyPrefix = []byte{0x00}
	AddressKeyPrefix = []byte{0x10}
)

func GetProfileKey(handle string) []byte {
	return append(ProfileKeyPrefix, []byte(handle)...)
}

func GetAddressKey(accAddr sdk.AccAddress) []byte {
	return append(AddressKeyPrefix, accAddr.Bytes()...)
}
