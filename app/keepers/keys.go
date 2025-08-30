package keepers

import (
	storetypes "cosmossdk.io/store/types"

	"cosmossdk.io/x/feegrant"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cadencetypes "github.com/bitsongofficial/go-bitsong/x/cadence/types"
	smartaccounttypes "github.com/bitsongofficial/go-bitsong/x/smart-account/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/types"

	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v10/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

func (appKeepers *AppKeepers) GenerateKeys() {
	appKeepers.keys = storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		crisistypes.StoreKey,
		paramstypes.StoreKey,
		consensusparamtypes.StoreKey,
		upgradetypes.StoreKey,
		feegrant.StoreKey,
		evidencetypes.StoreKey,
		packetforwardtypes.StoreKey,
		ibcexported.StoreKey,
		ibchookstypes.StoreKey,
		ibctransfertypes.StoreKey,
		ibcwasmtypes.StoreKey,
		capabilitytypes.StoreKey,
		authzkeeper.StoreKey,
		wasmtypes.StoreKey,
		fantokentypes.StoreKey,
		cadencetypes.StoreKey,
		smartaccounttypes.StoreKey,
		protocolpooltypes.StoreKey,
	)

	appKeepers.tkeys = storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	// memkeys are info stored  in RAM
	appKeepers.memKeys = storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
}

func (appKeepers *AppKeepers) GetKVStoreKey() map[string]*storetypes.KVStoreKey {
	return appKeepers.keys
}

func (appKeepers *AppKeepers) GetTransientStoreKey() map[string]*storetypes.TransientStoreKey {
	return appKeepers.tkeys
}

func (appKeepers *AppKeepers) GetMemoryStoreKey() map[string]*storetypes.MemoryStoreKey {
	return appKeepers.memKeys
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetKey(storeKey string) *storetypes.KVStoreKey {
	return appKeepers.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appKeepers *AppKeepers) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return appKeepers.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (appKeepers *AppKeepers) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return appKeepers.memKeys[storeKey]
}
