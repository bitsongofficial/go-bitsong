package types

import "github.com/ipfs/go-cid"

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
// - 0x00<cid_Bytes>: Track
// - 0x01<cid_Bytes>: Artist
var (
	KeyLastTrackID = []byte("lastTrackId")

	TrackKeyPrefix  = []byte{0x00}
	ArtistKeyPrefix = []byte{0x01}
)

func GetTrackKey(c string) []byte {
	cid, _ := cid.Decode(c)
	return append(TrackKeyPrefix, cid.Bytes()...)
}

func GetArtistKey(c string) []byte {
	cid, _ := cid.Decode(c)
	return append(ArtistKeyPrefix, cid.Bytes()...)
}
