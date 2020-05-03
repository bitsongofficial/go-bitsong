package types

const (
	// ModuleName is the name of the module
	ModuleName = "player"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	MinMonikerLength int = 4
	MaxMonikerLength int = 32
)

var (
	PlayersCountKey = []byte{0x00}
	PlayerKeyPrefix = []byte{0x01}
)

func PlayerKey(id PlayerID) []byte {
	return append(PlayerKeyPrefix, id.Bytes()...)
}
