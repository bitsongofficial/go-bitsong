package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "content"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	MaxNameLength int = 65
	MaxHashLength int = 1000
	MaxUriLength  int = 165
	MaxDaoLength  int = 15
)

// Keys for content store
// Items are stored with the following key: values
//
// - 0x00<uri_Bytes>: Content
var (
	ContentKeyPrefix = []byte{0x00}

	KeyLastHlsID = []byte("lastHlsId")
	HlsKeyPrefix = []byte{0x01}
)

func GetContentKey(uri string) []byte {
	return append(ContentKeyPrefix, []byte(uri)...)
}

// hls
func GetHlsKey(hlsID uint64) []byte {
	hlsIDBz := sdk.Uint64ToBigEndian(hlsID)
	return append(HlsKeyPrefix, hlsIDBz...)
}
