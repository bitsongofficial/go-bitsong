package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestKeeper_CreateAccount(t *testing.T) {
	ctx, keeper := SetupTestInput()
	pubKey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	handle := "test"

	account, err := keeper.CreateAccount(ctx, addr, pubKey, handle)
	require.NoError(t, err)
	require.Equal(t, account.Address, addr)
	require.Equal(t, account.Handle, handle)
}
