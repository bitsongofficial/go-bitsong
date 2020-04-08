package types

const (
	// ModuleName defines the desmos IBC name
	ModuleName = "desmosibc"

	// StoreKey is the store key string for desmos IBC
	StoreKey = ModuleName

	// RouterKey is the message route for desmos IBC
	RouterKey = ModuleName

	// QuerierRoute is the querier route for desmos IBC
	QuerierRoute = ModuleName

	// Version defines the current version the IBC desmos
	// module supports
	Version = "ics20-1"

	// PortID that desmos module binds to
	PortID = "desmos"
)
