package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	nftkeeper "github.com/bitsongofficial/go-bitsong/x/nft/keeper"

	storetypes "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/x/drop/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeKey     storetypes.StoreKey
	storeService store.KVStoreService
	logger       log.Logger

	nftKeeper *nftkeeper.Keeper

	Schema collections.Schema
	Drops  collections.Map[string, types.Drop]                           // (collectionDenom) -> Drop
	Rules  collections.Map[collections.Pair[string, string], types.Rule] // (collectionDenom, ruleId) -> Rule
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, storeService store.KVStoreService, nftKeeper *nftkeeper.Keeper, logger log.Logger) Keeper {
	logger = logger.With(log.ModuleKey, "x/"+types.ModuleName)

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeKey:     key,
		storeService: storeService,
		logger:       logger,
		nftKeeper:    nftKeeper,
		Drops: collections.NewMap(
			sb,
			types.DropsPrefix,
			"drops",
			collections.StringKey, codec.CollValue[types.Drop](cdc),
		),
		Rules: collections.NewMap(
			sb,
			types.RulesPrefix,
			"rules",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			codec.CollValue[types.Rule](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}
