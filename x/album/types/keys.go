package types

import (
	"encoding/binary"
)

const (
	ModuleName = "album"    // ModuleName is the name of the module
	StoreKey   = ModuleName // StoreKey is the store key string for album
	RouterKey  = ModuleName // RouterKey is the message route for album
)

// Keys for album store
// Items are stored with the following key: values
//
// - 0x00<albumID_Bytes>: Album
//
// - 0x01: nextAlbumID
//
// - 0x10<albumID_Bytes><trackID_Bytes>: Track
var (
	AlbumsKeyPrefix = []byte{0x00}
	AlbumIDKey      = []byte{0x01}

	TracksKeyPrefix = []byte{0x10}
)

// TracksKey gets the first part of the track key based on the albumID
func TracksKey(albumID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, albumID)
	return append(TracksKeyPrefix, bz...)
}

// TrackKey key of a specific track from the store
func TrackKey(albumID uint64, trackID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, trackID)

	return append(TracksKey(albumID), bz...)
}
