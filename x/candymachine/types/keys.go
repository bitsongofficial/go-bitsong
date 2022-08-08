package types

const (
	// ModuleName is the name of the module
	ModuleName = "candymachine"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey string = ModuleName
)

var (
	PrefixCandyMachine          = []byte{0x01}
	PrefixCandyMachineByEndTime = []byte{0x02}
)
