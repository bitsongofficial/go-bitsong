package types

import (
	"encoding/binary"
)

const (
	// ModuleName is the name of the module
	ModuleName = "track"

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

// Keys for track store
// Items are stored with the following key: values
//
// - 0x00<trackID_Bytes>: Track
var (
	KeyLastTrackID = []byte("lastTrackId")

	TrackKeyPrefix = []byte{0x00}
)

func GetTrackIDBytes(id uint64) []byte {
	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	return idBz
}

func GetTrackKey(id uint64) []byte {
	return append(TrackKeyPrefix, GetTrackIDBytes(id)...)
}
