package types

import (
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

	// TODO: move to params
	MaxTrackInfoLength = 30 * 1024
)

// Keys for track store
// Items are stored with the following key: values
//
// - 0x00<trackID_Bytes>: Track
// - 0x10<creatorAddr_Bytes><trackID_Bytes>: Track
var (
	TrackKeyPrefix         = []byte{0x00}
	TracksCreatorKeyPrefix = []byte{0x10}
)

func GetTrackIDBytes(id string) []byte {
	return []byte(id)
}

func GetTrackKey(id string) []byte {
	return append(TrackKeyPrefix, GetTrackIDBytes(id)...)
}

func GetCreatorKey(addr sdk.AccAddress) []byte {
	return append(TracksCreatorKeyPrefix, addr.Bytes()...)
}

func GetTrackByCreatorAddr(addr sdk.AccAddress, trackID string) []byte {
	return append(GetCreatorKey(addr), GetTrackIDBytes(trackID)...)
}
