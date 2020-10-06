package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func Test_handleMsgProfileCreate(t *testing.T) {
	ctx, _, keeper := SetupTestInput()
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	handle := "test"
	handle2 := "test2"
	metadataURI := "metadata"

	h := NewHandler(keeper)
	require.NotNil(t, h)

	msg := types.NewMsgChannelCreate(addr, handle, metadataURI)
	require.NotNil(t, msg)

	result, err := h(ctx, msg)
	require.NoError(t, err)

	account := types.Channel{}
	err = keeper.codec.UnmarshalBinaryLengthPrefixed(result.Data, &account)
	require.NoError(t, err)

	require.Equal(t, account.Handle, handle)
	require.Equal(t, account.Owner, addr)
	require.Equal(t, account.MetadataURI, metadataURI)

	msg = types.NewMsgChannelCreate(addr, handle, metadataURI)
	require.NotNil(t, msg)

	result, err = h(ctx, msg)
	require.Error(t, err)

	msg = types.NewMsgChannelCreate(addr, handle2, metadataURI)
	require.NotNil(t, msg)

	result, err = h(ctx, msg)
	require.Error(t, err)
}
