package types

const (
	// ModuleName is the name of the module
	ModuleName = "sdp"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

// Keys for sdp store
// Items are stored with the following key: values
//
// - 0x00<hash_Bytes>: Content
var (
	SdpKeyPrefix = []byte{0x00}
)

// TODO: to be refactored with a calculated hash
func GetSdpKey(hash string) []byte {
	return append(SdpKeyPrefix, []byte(hash)...)
}
