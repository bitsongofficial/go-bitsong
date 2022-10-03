package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "launchpad"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey string = ModuleName
)

var (
	PrefixLaunchPad           = []byte{0x01}
	PrefixLaunchPadByEndTime  = []byte{0x02}
	PrefixMintableMetadataIds = []byte{0x03}
)

func LaunchPadKey(collId uint64) []byte {
	return append(PrefixLaunchPad, sdk.Uint64ToBigEndian(collId)...)
}

func LaunchPadByEndTimeKey(timestamp uint64, collId uint64) []byte {
	idBz := sdk.Uint64ToBigEndian(collId)
	return append(GetTimeKey(timestamp), idBz...)
}

func CollectionMintableMetadataIdPrefix(collId uint64) []byte {
	return append(PrefixMintableMetadataIds, sdk.Uint64ToBigEndian(collId)...)
}

func GetTimeKey(timestamp uint64) []byte {
	time := time.Unix(int64(timestamp), 0)
	timeBz := sdk.FormatTimeBytes(time)
	timeBzL := len(timeBz)
	prefixL := len(PrefixLaunchPadByEndTime)

	bz := make([]byte, prefixL+8+timeBzL)

	// copy the prefix
	copy(bz[:prefixL], PrefixLaunchPadByEndTime)

	// copy the encoded time bytes length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))

	// copy the encoded time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)
	return bz
}
