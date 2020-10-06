package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "channel"

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
// - 0x00<handle_Bytes>: Channel
// - 0x10<accAddr_Bytes>: Handle
var (
	ChannelKeyPrefix = []byte{0x00}
	AddressKeyPrefix = []byte{0x10}
)

func GetChannelKey(handle string) []byte {
	return append(ChannelKeyPrefix, []byte(handle)...)
}

func GetOwnerKey(accAddr sdk.AccAddress) []byte {
	return append(AddressKeyPrefix, accAddr.Bytes()...)
}
