package types

const (
	// module name
	ModuleName = "track"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	DefaultStartingSongID uint64 = 1

	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName
)
