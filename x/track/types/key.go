package types

import (
	"github.com/tendermint/tendermint/crypto"
	"regexp"
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

	MaxTitleLength           int = 140
	MaxAttributesLength      int = 10
	MaxAttributesValueLength int = 50
)

// Keys for track store
// Items are stored with the following key: values
//
// - 0x00<trackAddr_Bytes>: Track
var (
	PathRegEx = regexp.MustCompile(`^\/(ip[fn]s)\/([^/?#]+)`)

	KeyLastTrackID = []byte("lastTrackId")
	TrackKeyPrefix = []byte{0x00}
)

func GetTrackKey(addr crypto.Address) []byte {
	return append(TrackKeyPrefix, addr...)
}
