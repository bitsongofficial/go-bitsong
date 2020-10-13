package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/auth/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestKeeper_RegisterHandle(t *testing.T) {
	ctx, _, ak, keeper := SetupTestInput()
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	_, err := keeper.RegisterHandle(ctx, addr, "handle")
	require.Error(t, err)

	acc := ak.NewAccountWithAddress(ctx, addr)
	err = acc.SetCoins(
		sdk.NewCoins(
			sdk.NewCoin("stake", sdk.OneInt()),
		),
	)
	ak.SetAccount(ctx, acc)
	require.NoError(t, err)

	bacc, err := keeper.RegisterHandle(ctx, addr, "handle")
	require.NoError(t, err)
	require.IsType(t, &types.BitSongAccount{}, bacc)
	require.Equal(t, bacc.GetAddress(), addr)
}
