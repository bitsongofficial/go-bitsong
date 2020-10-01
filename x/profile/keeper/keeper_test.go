package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestKeeper_CreateProfile(t *testing.T) {
	ctx, _, keeper := SetupTestInput()
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	handle := "test"
	handle2 := "test2"
	metadataURI := "metadata"

	account, err := keeper.CreateProfile(ctx, addr, handle, metadataURI)
	require.NoError(t, err)
	require.Equal(t, account.Address, addr)
	require.Equal(t, account.Handle, handle)
	require.Equal(t, account.MetadataURI, metadataURI)

	account, err = keeper.CreateProfile(ctx, addr, handle, metadataURI)
	require.Error(t, err)

	account, err = keeper.CreateProfile(ctx, addr, handle2, metadataURI)
	require.Error(t, err)
}
