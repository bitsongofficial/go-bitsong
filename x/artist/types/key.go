package types

import btsg "github.com/bitsongofficial/go-bitsong/types"

const (
	// ModuleName is the name of the module
	ModuleName = "artist"

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
// - 0x00<id_Bytes>: Artist
var (
	ArtistKeyPrefix = []byte{0x00}
)

func GetArtistKey(id btsg.ID) []byte {
	return append(ArtistKeyPrefix, []byte(id)...)
}
