package types

import (
	"encoding/binary"
	"fmt"
)

const (
	// ModuleName is the name of the module
	ModuleName = "artist"

	// StoreKey is the store key string for gov
	StoreKey = ModuleName

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

// Keys for artist store
// Items are stored with the following key: values
//
// - 0x00<artistID_Bytes>: Artist
//
// - 0x01: nextArtistID
var (
	ArtistsKeyPrefix = []byte{0x00}
	ArtistIDKey      = []byte{0x01}
)

// ArtistKey gets a specific artist from the store
func ArtistKey(artistID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, artistID)
	return append(ArtistsKeyPrefix, bz...)
}

// Split keys function; used for iterators

// SplitArtistKey split the artist key and returns the artist id
func SplitArtistKey(key []byte) (artistID uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d ≠ 8)", len(key[1:])))
	}

	return binary.LittleEndian.Uint64(key[1:])
}
