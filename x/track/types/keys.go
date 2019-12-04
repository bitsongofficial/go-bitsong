package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "track"    // ModuleName is the name of the module
	StoreKey   = ModuleName // StoreKey is the store key string for track
	RouterKey  = ModuleName // RouterKey is the message route for track
)

// Keys for track store
// Items are stored with the following key: values
//
// - 0x00<trackID_Bytes>: Track
//
// - 0x01: nextTrackID
//
// - 0x10<trackID_Bytes><accAddr_Bytes>: Play
var (
	TracksKeyPrefix = []byte{0x00}
	TrackIDKey      = []byte{0x01}

	PlaysKeyPrefix = []byte{0x10}
)

// PlaysKey gets the first part of the play key based on the trackID
func PlaysKey(trackId uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, trackId)
	return append(PlaysKeyPrefix, bz...)
}

// PlayKey key of a specific play from the store
func PlayKey(trackID uint64, accAddr sdk.AccAddress) []byte {
	return append(PlaysKey(trackID), accAddr.Bytes()...)
}
