package types

const (
	ModuleName = "album"    // ModuleName is the name of the module
	StoreKey   = ModuleName // StoreKey is the store key string for album
	RouterKey  = ModuleName // RouterKey is the message route for album
)

// Keys for album store
// Items are stored with the following key: values
//
// - 0x00<albumID_Bytes>: Album
//
// - 0x01: nextAlbumID
var (
	AlbumsKeyPrefix = []byte{0x00}
	AlbumIDKey      = []byte{0x01}
)
