package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/release/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestNewQuerier(t *testing.T) {
	ctx, _, keeper := SetupTestInput()
	querier := NewQuerier(keeper)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err := querier(ctx, []string{"foo", "bar"}, query)
	require.Error(t, err)
}

func Test_queryRelease(t *testing.T) {
	ctx, cdc, keeper := SetupTestInput()
	creator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	releaseID := "test1"
	metadataURI := "metadata"

	h := NewHandler(keeper)
	require.NotNil(t, h)

	msg := types.NewMsgReleseCreate(releaseID, metadataURI, creator)
	require.NotNil(t, msg)

	_, err := h(ctx, msg)
	require.NoError(t, err)

	querier := NewQuerier(keeper)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err = querier(ctx, []string{"foo", "bar"}, query)
	require.Error(t, err)

	query.Path = "/custom/release/release"
	var res []byte

	query.Data = []byte("?")
	res, err = querier(ctx, []string{types.QueryRelease}, query)
	require.Error(t, err)
	require.Nil(t, res)

	params := types.NewQueryReleaseParams(releaseID)
	bz, err := cdc.MarshalJSON(params)
	require.Nil(t, err)

	query.Data = bz
	res, err = querier(ctx, []string{types.QueryRelease}, query)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func Test_queryAllReleasesForCreator(t *testing.T) {
	ctx, cdc, keeper := SetupTestInput()
	creator := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	creator2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	releaseID := "test1"
	releaseID2 := "test2"
	metadataURI := "metadata"

	h := NewHandler(keeper)
	require.NotNil(t, h)

	msg := types.NewMsgReleseCreate(releaseID, metadataURI, creator)
	require.NotNil(t, msg)

	_, err := h(ctx, msg)
	require.NoError(t, err)

	querier := NewQuerier(keeper)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err = querier(ctx, []string{"foo", "bar"}, query)
	require.Error(t, err)

	query.Path = "/custom/release/release"
	var res []byte

	query.Data = []byte("?")
	res, err = querier(ctx, []string{types.QueryAllReleaseForCreator}, query)
	require.Error(t, err)
	require.Nil(t, res)

	params := types.NewQueryAllReleaseForCreatorParams(creator)
	bz, err := cdc.MarshalJSON(params)
	require.Nil(t, err)

	query.Data = bz
	res, err = querier(ctx, []string{types.QueryAllReleaseForCreator}, query)
	var releases []types.Release
	keeper.codec.MustUnmarshalJSON(res, &releases)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 1, len(releases))

	msg = types.NewMsgReleseCreate(releaseID2, metadataURI, creator)
	require.NotNil(t, msg)

	_, err = h(ctx, msg)
	require.NoError(t, err)

	res, err = querier(ctx, []string{types.QueryAllReleaseForCreator}, query)
	keeper.codec.MustUnmarshalJSON(res, &releases)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 2, len(releases))

	params = types.NewQueryAllReleaseForCreatorParams(creator2)
	bz, err = cdc.MarshalJSON(params)
	require.Nil(t, err)

	query.Data = bz

	res, err = querier(ctx, []string{types.QueryAllReleaseForCreator}, query)
	keeper.codec.MustUnmarshalJSON(res, &releases)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 0, len(releases))
}
