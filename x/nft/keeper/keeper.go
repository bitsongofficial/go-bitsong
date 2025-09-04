package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NFTIndexes struct {
	Collection *indexes.Multi[string, collections.Pair[string, string], types.Nft]
	Owner      *indexes.Multi[sdk.AccAddress, collections.Pair[string, string], types.Nft]
}

func newNFTIndexes(sb *collections.SchemaBuilder) NFTIndexes {
	return NFTIndexes{
		Collection: indexes.NewMulti(
			sb,
			types.NFTsByCollectionPrefix,
			"nfts_by_collection",
			collections.StringKey,
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			func(pk collections.Pair[string, string], v types.Nft) (string, error) {
				return v.Collection, nil
			},
		),
		Owner: indexes.NewMulti(
			sb,
			types.OwnersPrefix,
			"owners",
			sdk.AccAddressKey,
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			func(pk collections.Pair[string, string], v types.Nft) (sdk.AccAddress, error) {
				return sdk.AccAddressFromBech32(v.Owner)
			},
		),
	}
}

func (i NFTIndexes) IndexesList() []collections.Index[collections.Pair[string, string], types.Nft] {
	return []collections.Index[collections.Pair[string, string], types.Nft]{
		i.Collection,
		i.Owner,
	}
}

type Keeper struct {
	cdc          codec.BinaryCodec
	storeKey     storetypes.StoreKey
	storeService store.KVStoreService
	ac           address.Codec
	// bk           types.BankKeeper
	logger log.Logger

	Schema      collections.Schema
	Collections collections.Map[string, types.Collection]
	Supply      collections.Map[string, math.Int]
	// (collectionDenom, tokenId) -> NFT
	NFTs *collections.IndexedMap[collections.Pair[string, string], types.Nft, NFTIndexes]
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, storeService store.KVStoreService, ak types.AccountKeeper, logger log.Logger) Keeper {
	/*if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the " + types.ModuleName + " module account has not been set")
	}*/

	logger = logger.With(log.ModuleKey, "x/"+types.ModuleName)

	sb := collections.NewSchemaBuilder(storeService)
	ac := ak.AddressCodec()

	k := Keeper{
		cdc:          cdc,
		storeKey:     key,
		storeService: storeService,
		ac:           ac,
		logger:       logger,
		// TODO: fix the store once we add queries
		Collections: collections.NewMap(sb, types.CollectionsPrefix, "collections", collections.StringKey, codec.CollValue[types.Collection](cdc)),
		Supply:      collections.NewMap(sb, types.SupplyPrefix, "supply", collections.StringKey, sdk.IntValue),
		NFTs: collections.NewIndexedMap(
			sb,
			types.NFTsPrefix,
			"nfts",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			codec.CollValue[types.Nft](cdc),
			newNFTIndexes(sb),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}
