package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/bitsong/x/fantoken/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryToken() {
	app, ctx := suite.app, suite.ctx
	_, _, addr := testdata.KeyTestPubAddr()
	token := types.NewFanToken("btc", "Bitcoin Token", sdk.NewInt(22000000), true, "test", addr)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.FanTokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_ = suite.app.FanTokenKeeper.AddFanToken(ctx, token)

	// Query token
	tokenResp1, err := queryClient.FanToken(gocontext.Background(), &types.QueryFanTokenRequest{Denom: "btc"})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokenResp1)

	tokenResp2, err := queryClient.FanToken(gocontext.Background(), &types.QueryFanTokenRequest{Denom: "ubtc"})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokenResp2)

	// Query tokens
	tokensResp1, err := queryClient.FanTokens(gocontext.Background(), &types.QueryFanTokensRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokensResp1)
	suite.Len(tokensResp1.FanTokens, 1)
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
	token := types.NewFanToken("btc", "Bitcoin Token", sdk.NewInt(22000000), true, "test", addr)
	err := suite.app.FanTokenKeeper.AddFanToken(ctx, token)
	suite.Require().NoError(err)

	buinCoin := sdk.NewInt64Coin("satoshi", 1000000000000000000)
	app.FanTokenKeeper.AddBurnCoin(ctx, buinCoin)

	resp, err := queryClient.TotalBurn(gocontext.Background(), &types.QueryTotalBurnRequest{})
	suite.Require().NoError(err)
	suite.Len(resp.BurnedCoins, 1)
	suite.EqualValues(buinCoin, resp.BurnedCoins[0])
}
