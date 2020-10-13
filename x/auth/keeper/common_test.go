package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/auth/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func SetupTestInput() (sdk.Context, *codec.Codec, auth.AccountKeeper, Keeper) {
	// define store keys
	accKey := sdk.NewKVStoreKey(types.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	tKey := sdk.NewTransientStoreKey(params.TStoreKey)

	// create an in-memory db
	memDB := db.NewMemDB()

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tKey, sdk.StoreTypeIAVL, memDB)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	// create a Cdc and a context
	cdc := testCodec()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	paramsKeeper := params.NewKeeper(cdc, paramsKey, tKey)
	ak := auth.NewAccountKeeper(cdc, accKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bak := NewKeeper(ak)

	return ctx, cdc, ak, bak
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)

	codec.RegisterCrypto(cdc)

	cdc.Seal()
	return cdc
}
