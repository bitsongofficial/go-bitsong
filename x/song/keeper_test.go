package song

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"testing"
)

type TestInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	k   Keeper
}

func SetupTestInput() TestInput {
	cdc := codec.New()

	songCapKey := sdk.NewKVStoreKey("songCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(songCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	k := NewKeeper(songCapKey, cdc)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-ID"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	return TestInput{cdc: cdc, ctx: ctx, k: k}
}

func TestKeeper(t *testing.T) {
	input := SetupTestInput()
	ctx := input.ctx
	k := input.k

	_, err := k.Publish(ctx, "Test Song", sdk.AccAddress([]byte("addr1")), "", sdk.NewDecWithPrec(5, 2))
	require.NoError(t, err)
}
