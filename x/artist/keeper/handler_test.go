package keeper

import (
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func Test_handleMsgArtistCreate(t *testing.T) {
	ctx, _, keeper := SetupTestInput()
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	id := btsg.ID("dasdas-dasdsa")
	name := "Freedy Mercury"

	h := NewHandler(keeper)
	require.NotNil(t, h)

	msg := types.NewMsgArtistCreate(id.String(), name, nil, nil, "", addr)
	require.NotNil(t, msg)

	result, err := h(ctx, msg)
	require.NoError(t, err)

	artist := types.Artist{}
	err = keeper.codec.UnmarshalBinaryLengthPrefixed(result.Data, &artist)
	require.NoError(t, err)

	require.Equal(t, artist.ID, id)
	require.Equal(t, artist.Name, name)
	require.Equal(t, artist.Creator, addr)

	msg = types.NewMsgArtistCreate(id.String(), name, nil, nil, "", addr)
	require.NotNil(t, msg)

	result, err = h(ctx, msg)
	require.Error(t, err)
}
