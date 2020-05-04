package types

const (
	// ModuleName is the name of the module
	ModuleName = "content"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	MaxNameLength          int = 65
	MaxMetadataLength      int = 1000
	MaxUriLength           int = 165
	MaxRightsHoldersLength int = 15
)

// Keys for content store
// Items are stored with the following key: values
//
// - 0x00<uri_Bytes>: Content
// - 0x01<denom_Bytes>: Denom
var (
	ContentKeyPrefix = []byte{0x00}
	DenomKeyPrefix   = []byte{0x01}
)

func GetContentKey(uri string) []byte {
	return append(ContentKeyPrefix, []byte(uri)...)
}

func GetDenomKey(denom string) []byte {
	return append(DenomKeyPrefix, []byte(denom)...)
}
