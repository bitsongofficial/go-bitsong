package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/bitsongofficial/chainmodules/x/fantoken/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryToken() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ubtc",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ubtc", Exponent: 0},
			{Denom: "btc", Exponent: types.FanTokenDecimal},
		},
	}
	token := types.NewFanToken("Bitcoin Network", sdk.NewInt(22000000), addr, denomMetaData)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.FanTokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_ = suite.app.FanTokenKeeper.AddFanToken(ctx, token)

	// Query token
	tokenResp, err := queryClient.FanToken(gocontext.Background(), &types.QueryFanTokenRequest{Denom: "ubtc"})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokenResp)

	suite.Require().Equal(tokenResp.Token.Name, "Bitcoin Network")
	suite.Require().Equal(tokenResp.Token.MaxSupply, sdk.NewInt(22000000))
	suite.Require().Equal(tokenResp.Token.Owner, addr.String())
	suite.Require().Equal(tokenResp.Token.MetaData, denomMetaData)

	// Query tokens
	tokensResp, err := queryClient.FanTokens(gocontext.Background(), &types.QueryFanTokensRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokensResp)
	suite.Len(tokensResp.Tokens, 1)
}

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.FanTokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	params := app.FanTokenKeeper.GetParamSet(ctx)
	suite.Require().NoError(err)
	suite.Equal(params, paramsResp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryTotalBurn() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.FanTokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, _, addr := testdata.KeyTestPubAddr()
	denomMetaData := banktypes.Metadata{
		Description: "test",
		Base:        "ubtc",
		Display:     "btc",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ubtc", Exponent: 0},
			{Denom: "btc", Exponent: types.FanTokenDecimal},
		},
	}
	token := types.NewFanToken("Bitcoin Network", sdk.NewInt(22000000), addr, denomMetaData)
	err := suite.app.FanTokenKeeper.AddFanToken(ctx, token)
	suite.Require().NoError(err)

	burnCoin := sdk.NewInt64Coin("satoshi", 1000000000000000000)
	app.FanTokenKeeper.AddBurnCoin(ctx, burnCoin)

	resp, err := queryClient.TotalBurn(gocontext.Background(), &types.QueryTotalBurnRequest{})
	suite.Require().NoError(err)
	suite.Len(resp.BurnedCoins, 1)
	suite.EqualValues(burnCoin, resp.BurnedCoins[0])
}
