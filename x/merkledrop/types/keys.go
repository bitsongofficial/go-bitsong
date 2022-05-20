package types

const (
	// ModuleName is the name of the module
	ModuleName = "merkledrop"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey string = ModuleName
)

var (
	PrefixMerkleDrop        = []byte{0x01}
	PrefixMerkleDropByOwner = []byte{0x02}
	KeyLastMerkleDropId     = []byte{0x03}
)
