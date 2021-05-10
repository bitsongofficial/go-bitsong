package fantoken_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/bitsongofficial/bitsong/app"
	token "github.com/bitsongofficial/bitsong/x/fantoken"
	"github.com/bitsongofficial/bitsong/x/fantoken/types"
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// export genesis
	genesisState := token.ExportGenesis(ctx, app.FanTokenKeeper)

	require.Equal(t, types.DefaultParams(), genesisState.Params)
	for _, token := range genesisState.Tokens {
		require.Equal(t, token, types.GetNativeToken())
	}
}

func TestInitGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// add token
	addr := sdk.AccAddress(tmhash.SumTruncated([]byte("addr1")))
	ft := types.NewFanToken("btc", "Bitcoin Network", sdk.NewInt(1), true, addr)

	burnCoins := []sdk.Coin{
		{Denom: ft.Denom, Amount: sdk.NewInt(1000)},
	}
	genesis := types.GenesisState{
		Params:      types.DefaultParams(),
		Tokens:      []types.FanToken{ft},
		BurnedCoins: burnCoins,
	}

	// initialize genesis
	token.InitGenesis(ctx, app.FanTokenKeeper, genesis)

	// query all tokens
	var tokens = app.FanTokenKeeper.GetFanTokens(ctx, nil)
	require.Equal(t, len(tokens), 2)
	require.Equal(t, tokens[0], &ft)

	var coins = app.FanTokenKeeper.GetAllBurnCoin(ctx)
	require.Equal(t, burnCoins, coins)
}
