package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func SetupTestInput() (sdk.Context, *codec.Codec, Keeper) {
	// define store keys
	artistKey := sdk.NewKVStoreKey(types.StoreKey)
	//paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	//tKey := sdk.NewTransientStoreKey(params.TStoreKey)

	// create an in-memory db
	memDB := db.NewMemDB()

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(artistKey, sdk.StoreTypeIAVL, memDB)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	// create a Cdc and a context
	cdc := testCodec()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	//paramsKeeper := params.NewKeeper(cdc, paramsKey, tKey)
	artistKeeper := NewKeeper(artistKey, cdc)

	return ctx, cdc, artistKeeper
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	types.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}
