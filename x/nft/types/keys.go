package types

const (
	// ModuleName is the name of the module
	ModuleName = "nft"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey string = ModuleName
)

var (
	PrefixNFT             = []byte{0x01}
	PrefixMetadata        = []byte{0x02}
	PefixCollection       = []byte{0x03}
	PefixCollectionRecord = []byte{0x04}
)
