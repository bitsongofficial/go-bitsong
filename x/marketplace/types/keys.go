package types

const (
	// ModuleName is the name of the module
	ModuleName = "marketplace"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey string = ModuleName
)

var (
	PrefixAuction            = []byte{0x01}
	PrefixAuctionByAuthority = []byte{0x02}
	PrefixAuctionByEndTime   = []byte{0x03}
	KeyLastAuctionId         = []byte{0x04}
	PrefixBid                = []byte{0x05}
	PrefixBidByBidder        = []byte{0x06}
	PrefixBidderMetadata     = []byte{0x07}
)
