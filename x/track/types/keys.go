package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "track"    // ModuleName is the name of the module
	StoreKey   = ModuleName // StoreKey is the store key string for track
	RouterKey  = ModuleName // RouterKey is the message route for track

	DefaultParamspace = ModuleName // DefaultParamspace default name for parameter store
)

// Keys for track store
// Items are stored with the following key: values
//
// - 0x00<trackID_Bytes>: Track
//
// - 0x01: nextTrackID
//
// - 0x10<trackID_Bytes><accAddr_Bytes>: Play
//
// - 0x20<trackID_Bytes>: Share
//
// - 0x30<trackID_Bytes><accAddr_Bytes>: Deposit
var (
	TracksKeyPrefix = []byte{0x00}
	TrackIDKey      = []byte{0x01}

	PlaysKeyPrefix    = []byte{0x10}
	SharesKeyPrefix   = []byte{0x20}
	DepositsKeyPrefix = []byte{0x30}
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

func ShareKey(trackID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, trackID)
	return append(SharesKeyPrefix, bz...)
}

func DepositsKey(trackID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, trackID)
	return append(DepositsKeyPrefix, bz...)
}

// DepositKey key of a specific deposit from the store
func DepositKey(trackID uint64, depositorAddr sdk.AccAddress) []byte {
	return append(DepositsKey(trackID), depositorAddr.Bytes()...)
}
