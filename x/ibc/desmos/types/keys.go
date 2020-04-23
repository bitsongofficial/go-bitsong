package types

const (
	// ModuleName defines the IBC desmos name
	ModuleName = "desmosibc"

	// Version defines the current version the IBC desmos
	// module supports
	Version = "ics20-1"

	// PortID that IBC desmos module binds to
	PortID = "desmos"

	// StoreKey is the store key string for IBC desmos
	StoreKey = ModuleName

	// RouterKey is the message route for IBC desmos
	RouterKey = ModuleName

	// Key to store portID in our store
	PortKey = "portID"

	// QuerierRoute is the querier route for IBC desmos
	QuerierRoute = ModuleName
)
