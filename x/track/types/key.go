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
// - 0x20<trackID_Bytes><entityAddr_Bytes>: Share
var (
	TrackKeyPrefix         = []byte{0x00}
	TracksCreatorKeyPrefix = []byte{0x10}
	SharesKeyPrefix        = []byte{0x20}
)

func GetTrackIDBytes(trackID string) []byte {
	return []byte(trackID)
}

func GetTrackKey(trackID string) []byte {
	return append(TrackKeyPrefix, GetTrackIDBytes(trackID)...)
}

func GetCreatorKey(addr sdk.AccAddress) []byte {
	return append(TracksCreatorKeyPrefix, addr.Bytes()...)
}

func GetTrackByCreatorAddr(addr sdk.AccAddress, trackID string) []byte {
	return append(GetCreatorKey(addr), GetTrackIDBytes(trackID)...)
}

func SharesKey(trackID string) []byte {
	return append(SharesKeyPrefix, GetTrackIDBytes(trackID)...)
}

func ShareKey(trackID string, entityAddr sdk.AccAddress) []byte {
	return append(SharesKey(trackID), entityAddr.Bytes()...)
}

/*func GetTrackShareKey(id string) []byte {
	return append(TracksDepositKeyPrefix, GetTrackIDBytes(id)...)
}

func GetSharesByTrackIDAndEntity(id string, entity sdk.AccAddress) []byte {
	return append(GetTrackShareKey(id), entity.Bytes()...)
}*/
