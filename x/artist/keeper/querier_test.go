package keeper

import (
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
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

func Test_queryArtist(t *testing.T) {
	ctx, cdc, keeper := SetupTestInput()
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	id := btsg.ID("dasdas-dasdsa")
	name := "Freedy Mercury"

	h := NewHandler(keeper)
	require.NotNil(t, h)

	msg := types.NewMsgArtistCreate(id.String(), name, nil, nil, "", addr)
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

	query.Path = "/custom/artist/artist"
	var res []byte

	query.Data = []byte("?")
	res, err = querier(ctx, []string{types.QueryArtist}, query)
	require.Error(t, err)
	require.Nil(t, res)

	params := types.NewQueryArtistParams(id.String())
	bz, err := cdc.MarshalJSON(params)
	require.Nil(t, err)

	query.Data = bz
	res, err = querier(ctx, []string{types.QueryArtist}, query)
	require.NoError(t, err)
	require.NotNil(t, res)
}
