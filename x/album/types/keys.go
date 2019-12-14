package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName        = "album"    // ModuleName is the name of the module
	StoreKey          = ModuleName // StoreKey is the store key string for album
	RouterKey         = ModuleName // RouterKey is the message route for album
	DefaultParamspace = ModuleName // DefaultParamspace default name for parameter store
)

// Keys for album store
// Items are stored with the following key: values
//
// - 0x00<albumID_Bytes>: Album
//
// - 0x01: nextAlbumID
//
// - 0x10<albumID_Bytes><trackID_Bytes>: Track
//
// - 0x20<artistID_Bytes><depositorAddr_Bytes>: Deposit
var (
	AlbumsKeyPrefix = []byte{0x00}
	AlbumIDKey      = []byte{0x01}

	TracksKeyPrefix   = []byte{0x10}
	DepositsKeyPrefix = []byte{0x20}
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

// DepositsKey gets the first part of the deposits key based on the albumID
func DepositsKey(albumID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, albumID)
	return append(DepositsKeyPrefix, bz...)
}

// DepositKey key of a specific deposit from the store
func DepositKey(albumID uint64, depositorAddr sdk.AccAddress) []byte {
	return append(DepositsKey(albumID), depositorAddr.Bytes()...)
}
