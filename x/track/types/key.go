package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	MaxTrackInfoLength = 2 * 2024
)

// Keys for track store
// Items are stored with the following key: values
//
// - 0x00<trackID_Bytes>: Track
// - 0x10<creatorAddr_Bytes><trackID_Bytes>: Track
var (
	KeyLastTrackID = []byte("lastTrackId")

	TrackKeyPrefix         = []byte{0x00}
	TracksCreatorKeyPrefix = []byte{0x10}
)

func GetTrackIDBytes(id uint64) []byte {
	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	return idBz
}

func GetTrackKey(id uint64) []byte {
	return append(TrackKeyPrefix, GetTrackIDBytes(id)...)
}

func GetCreatorKey(addr sdk.AccAddress) []byte {
	return append(TracksCreatorKeyPrefix, addr.Bytes()...)
}

func GetTrackByCreatorAddr(addr sdk.AccAddress, trackID uint64) []byte {
	return append(GetCreatorKey(addr), GetTrackIDBytes(trackID)...)
}
