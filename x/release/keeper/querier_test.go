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

/*func Test_queryProfileByAddress(t *testing.T) {
	ctx, cdc, keeper := SetupTestInput()
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	handle := "test"
	metadataURI := "metadata"

	h := NewHandler(keeper)
	require.NotNil(t, h)

	msg := types.NewMsgProfileCreate(addr, handle, metadataURI)
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

	query.Path = "/custom/profile/profile"
	var res []byte

	query.Data = []byte("?")
	res, err = querier(ctx, []string{types.QueryProfileByAddress}, query)
	require.Error(t, err)
	require.Nil(t, res)

	params := types.NewQueryByAddressParams(addr)
	bz, err := cdc.MarshalJSON(params)
	require.Nil(t, err)

	query.Data = bz
	res, err = querier(ctx, []string{types.QueryProfileByAddress}, query)
	require.NoError(t, err)
	require.NotNil(t, res)
}
*/
