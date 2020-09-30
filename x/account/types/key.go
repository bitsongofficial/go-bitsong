package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "account"

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
// - 0x00<accAddr_Bytes>: Account
var (
	AccountKeyPrefix = []byte{0x00}
)

func GetAccountKey(accAddr sdk.AccAddress) []byte {
	return append(AccountKeyPrefix, accAddr.Bytes()...)
}
