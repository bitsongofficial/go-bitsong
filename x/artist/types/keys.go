package types

const (
	ModuleName = "artist"   // ModuleName is the name of the module
	StoreKey   = ModuleName // StoreKey is the store key string for artist
	RouterKey  = ModuleName // RouterKey is the message route for artist
)

// Keys for artist store
// Items are stored with the following key: values
//
// - 0x00<artistID_Bytes>: Artist
//
// - 0x01: nextArtistID
var (
	ArtistsKeyPrefix = []byte{0x00}
	ArtistIDKey      = []byte{0x01}
)
