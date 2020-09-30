package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/account/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func SetupTestInput() (sdk.Context, Keeper) {
	// define store keys
	authKey := sdk.NewKVStoreKey(types.StoreKey)
	accountKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	tKey := sdk.NewTransientStoreKey(params.TStoreKey)

	// create an in-memory db
	memDB := db.NewMemDB()

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, memDB)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	// create a Cdc and a context
	cdc := testCodec()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	paramsKeeper := params.NewKeeper(cdc, paramsKey, tKey)
	accountKeeper := auth.NewAccountKeeper(cdc, accountKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	authKeeper := NewKeeper(authKey, cdc, accountKeeper)

	return ctx, authKeeper
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)
	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	types.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}
