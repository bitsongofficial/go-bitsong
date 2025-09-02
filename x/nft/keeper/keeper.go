package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	ac           address.Codec
	bk           types.BankKeeper
	logger       log.Logger

	Schema      collections.Schema
	Collections collections.Map[string, types.Collection]
	Supply      collections.Map[string, math.Int]
}

func NewKeeper(cdc codec.BinaryCodec, storeService store.KVStoreService, ak types.AccountKeeper, bk types.BankKeeper, logger log.Logger) Keeper {
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the " + types.ModuleName + " module account has not been set")
	}

	logger = logger.With(log.ModuleKey, "x/"+types.ModuleName)

	sb := collections.NewSchemaBuilder(storeService)

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		ac:           ak.AddressCodec(),
		bk:           bk,
		logger:       logger,
		// TODO: fix the store once we add queries
		Collections: collections.NewMap(sb, types.CollectionsPrefix, "collections", collections.StringKey, codec.CollValue[types.Collection](cdc)),
		Supply:      collections.NewMap(sb, types.SupplyPrefix, "supply", collections.StringKey, sdk.IntValue),
	}
}
