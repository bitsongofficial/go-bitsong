package fantoken_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	token "github.com/bitsongofficial/go-bitsong/x/fantoken"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

var (
	owner     = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	name      = "Bitcoin"
	symbol    = "btc"
	uri       = "ipfs://"
	maxSupply = sdk.NewInt(200000000)
	mintable  = true
	initAmt   = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin  = sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initAmt)}
	issueFee  = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000))
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// export genesis
	genesisState := token.ExportGenesis(ctx, app.FanTokenKeeper)

	require.Equal(t, types.DefaultParams(), genesisState.Params)
}

func TestInitGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// add token
	addr := sdk.AccAddress(tmhash.SumTruncated([]byte("addr1")))
	fantokenObj := types.NewFanToken(name, symbol, uri, maxSupply, addr)

	burnCoins := []sdk.Coin{
		{Denom: fantokenObj.GetDenom(), Amount: sdk.NewInt(1000)},
	}
	genesis := types.GenesisState{
		Params:      types.DefaultParams(),
		FanTokens:   []types.FanToken{fantokenObj},
		BurnedCoins: burnCoins,
	}

	// initialize genesis
	token.InitGenesis(ctx, app.FanTokenKeeper, genesis)

	// query all fantokens
	var fantokens = app.FanTokenKeeper.GetFanTokens(ctx, nil)
	require.Equal(t, len(fantokens), 1)
	require.Equal(t, fantokens[0], &fantokenObj)

	var coins = app.FanTokenKeeper.GetAllBurnCoin(ctx)
	require.Equal(t, burnCoins, coins)
}
