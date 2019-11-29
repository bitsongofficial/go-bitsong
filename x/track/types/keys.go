package types

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
var (
	TracksKeyPrefix = []byte{0x00}
	TrackIDKey      = []byte{0x01}
)
